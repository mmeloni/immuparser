package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/codenotary/immudb/pkg/client"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

func Parse() error {

	if !IsLower(Ledger) {
		log.Fatalf("only lower case ledger names are allowed: %s", Ledger)
	}

	f, err := os.Open(Source)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	c, err := os.Open(Source)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	linesNumber, err := LineCounter(c)
	if err != nil {
		log.Fatal(err)
	}

	if linesNumber <= FlushSize {
		log.Fatalf("Flush size is greater than lines number: %d. Can not continue", linesNumber)
	}

	var keys []io.Reader
	var values []io.Reader

	cli, err := client.NewImmuClient(client.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	lr, err := cli.Login(ctx, []byte(`immudb`), []byte(`immudb`))
	if err != nil {
		log.Fatal(err)
	}
	md := metadata.Pairs("authorization", lr.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	err = cli.CreateDatabase(ctx, &schema.Database{
		Databasename: Ledger,
	})
	if err != nil {
		st, _ := status.FromError(err)
		if strings.Contains(st.Message(), "already exists") {
			log.Fatalf("ledger %s already present. Can not continue", Ledger)
		} else {
			log.Fatal(err)
		}
	}

	respUse, err := cli.UseDatabase(ctx, &schema.Database{
		Databasename: Ledger,
	})
	if err != nil {
		log.Fatal(err)
	}
	md = metadata.Pairs("authorization", respUse.Token)
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	scanner := bufio.NewScanner(f)
	count := 0
	gcount := 0
	bar := progressbar.Default(int64(linesNumber))

	for scanner.Scan() {
		keys = append(keys, strings.NewReader(fmt.Sprintf("row-%d", gcount)))
		values = append(values, strings.NewReader(strings.Replace(scanner.Text(), " ", "|", -1)))
		count++
		gcount++

		if err := bar.Add(1); err != nil {
			log.Fatal(err)
		}

		if count >= FlushSize {
			br := &client.BatchRequest{
				Keys:   keys,
				Values: values,
			}
			_, err = cli.SetBatch(ctx, br)
			if err != nil {
				log.Fatal(err)
			}
			keys = nil
			values = nil
			count = 0
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading "+Source+":", err)
	}
	return nil
}

func LineCounter(r io.Reader) (int, error) {

	var count int
	const lineBreak = '\n'

	buf := make([]byte, bufio.MaxScanTokenSize)

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			count++
		}
		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
