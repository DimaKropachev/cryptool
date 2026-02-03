package app

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/DimaKropachev/cryptool/pkg/crypto"
	"github.com/DimaKropachev/cryptool/pkg/crypto/algorithms"
	"github.com/DimaKropachev/cryptool/pkg/file"
	mem "github.com/DimaKropachev/cryptool/pkg/memory"
	"github.com/DimaKropachev/cryptool/pkg/table"
)

const (
	defaultIterations = 20
)

type ShortB struct {
	Alg            string
	Time           string
	MemUsed        string
	ResultFileSize string
}

func Benchmark(inPath string) error {
	inFileInfo, err := os.Stat(inPath)
	if err != nil {
		return fmt.Errorf("error get file info: %w", err)
	}

	info := `FILE INFO
Name: %s
Path: %s
Size: %s

`

	inFileSize := inFileInfo.Size()

	fmt.Printf(info,
		inFileInfo.Name(),
		inPath,
		mem.FormatBytes(float64(inFileSize)),
	)

	algs := []string{algorithms.AlgAES128GCM, algorithms.AlgAES192GCM, algorithms.AlgAES256GCM, algorithms.AlgCHACHA20POLY1305}
	result := make([][]string, len(algs))

	for i, alg := range algs {
		runtime.GC()

		var avgTime time.Duration
		var avgMem float64
		var outFileSize int64
		for j := 1; j < defaultIterations; j++ {

			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)

			start := time.Now()

			outFileSize, err = encrypt(inPath, int(inFileSize), alg)
			if err != nil {
				fmt.Println(err)
			}

			totalTime := time.Since(start)

			runtime.ReadMemStats(&after)

			avgTime += totalTime
			avgMem += float64(after.TotalAlloc - before.TotalAlloc)
		}
		avgTime /= defaultIterations
		avgMem /= defaultIterations
		result[i] = []string{alg, mem.FormatTime(avgTime), mem.FormatBytes(avgMem), mem.FormatBytes(float64(outFileSize))}
	}
	table := table.New()
	headlines := []string{"Algorithm", "Time", "Memory usage", "Result Size"}
	table.SetHeader(headlines)
	err = table.SetContent(result)
	if err != nil {
		return err
	}
	err = table.Render()
	if err != nil {
		return err
	}

	return nil
}

func encrypt(inPath string, inSize int, algName string) (int64, error) {
	var (
		salt = crypto.GenerateSalt(16)
	)

	alg, id, err := algorithms.CreateAlgorithmByName(algName, nil, salt)
	if err != nil {
		return 0, err
	}

	out, err := os.CreateTemp("", "bench.*.crpt")
	if err != nil {
		return 0, err
	}
	defer os.Remove(out.Name())
	defer out.Close()

	blockSize, err := CalculateOptimalBlockSize(inSize)
	if err != nil {
		return 0, err
	}

	header := crypto.NewHeader(id, blockSize, len(salt), alg.GetNonceSize(), salt)
	encHeader, err := crypto.EncryptHeader(header)
	if err != nil {
		return 0, fmt.Errorf("error encrypting header: %w", err)

	}

	if _, err := out.Write(encHeader); err != nil {
		return 0, err
	}

	content, errs, err := file.ReadDecryptedFile(inPath, blockSize)
	if err != nil {
		return 0, err
	}

READ:
	for {
		select {
		case plaintext, ok := <-content:
			if !ok {
				break READ
			}

			ciphertext, err := alg.Encrypt(plaintext)
			if err != nil {
				return 0, err
			}

			if _, err = out.Write(ciphertext); err != nil {
				return 0, err
			}

		case err, ok := <-errs:
			if !ok {
				break READ
			}
			return 0, fmt.Errorf("error reading file: %w", err)
		}
	}

	outFileInfo, err := out.Stat()
	if err != nil {
		return 0, fmt.Errorf("error get file info: %w", err)
	}

	return outFileInfo.Size(), nil
}
