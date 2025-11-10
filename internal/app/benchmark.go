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
	operationEncrypt = "encrypt"
	operationDecrypt = "decrypt"

	defaultIterations = 10
)

type ShortB struct {
	Alg     string
	Time    string
	MemUsed string
}

func Benchmark(inputPath string) error {
	err := file.ValidateFilePath(inputPath)
	if err != nil {
		return err
	}

	in, err := os.OpenFile(inputPath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("error to opening file \"%s\": %w", inputPath, err)
	}

	var operation string

	_, err = crypto.DecryptHeader(in)
	if err != nil {
		fmt.Println(err)
		operation = operationEncrypt
	} else {
		operation = operationDecrypt
	}

	switch operation {
	case operationDecrypt:

	case operationEncrypt:
		algs := []string{algorithms.AlgAES128GCM, algorithms.AlgAES192GCM, algorithms.AlgAES256GCM, algorithms.AlgCHACHA20POLY1305}
		// results := make([]ShortB, len(algs))
		result := make([][]string, len(algs))

		for i, alg := range algs {
			runtime.GC()

			var avgTime time.Duration
			var avgMem float64
			for i := 1; i < defaultIterations; i++ {

				var before, after runtime.MemStats
				runtime.ReadMemStats(&before)

				start := time.Now()

				encrypt(alg, inputPath)

				totalTime := time.Since(start)

				runtime.ReadMemStats(&after)

				avgTime += totalTime
				avgMem += float64(after.TotalAlloc - before.TotalAlloc)
			}
			avgTime /= defaultIterations
			avgMem /= defaultIterations
			// results[i] = ShortB{
			// 	Alg:     alg,
			// 	Time:    mem.FormatTime(avgTime),
			// 	MemUsed: mem.FormatBytes(avgMem),
			// }
			result[i] = []string{alg, mem.FormatTime(avgTime), mem.FormatBytes(avgMem)}
		}
		table := table.New()
		headlines := []string{"Algorithm", "Time", "Memory usage"}
		table.SetHeader(headlines)
		err := table.SetContent(result)
		if err != nil {
			return err
		}
		err = table.Render()
		if err != nil {
			return err
		}
	}

	return nil
}

func encrypt(algName string, inputPath string) error {
	var (
		salt = crypto.GenerateSalt(16)
	)

	in, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("error receiving information about an input file: %w", err)
	}

	alg, id, err := algorithms.CreateAlgorithmByName(algName, nil, salt)
	if err != nil {
		return err
	}

	out, err := os.CreateTemp("", "bench.*.crpt")
	if err != nil {
		return err
	}
	defer os.Remove(out.Name())
	defer out.Close()

	blockSize, err := CalculateOptimalBlockSize(int(in.Size()))
	if err != nil {
		return err
	}

	header := crypto.NewHeader(id, blockSize, len(salt), alg.GetNonceSize(), salt)
	encHeader, err := crypto.EncryptHeader(header)
	if err != nil {
		return fmt.Errorf("error encrypting header: %w", err)

	}

	if _, err := out.Write(encHeader); err != nil {
		return err
	}

	content, errs, err := file.ReadDecryptedFile(inputPath, blockSize)
	if err != nil {
		return err
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
				return err
			}

			if _, err = out.Write(ciphertext); err != nil {
				return err
			}

		case err, ok := <-errs:
			if !ok {
				break READ
			}
			return fmt.Errorf("error reading file: %w", err)
		}
	}

	return nil
}
