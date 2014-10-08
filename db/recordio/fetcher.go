package recordio

import "io"

type Fetcher interface {
	Fetch(offset int64, d Decoder) error
}

type fetcher struct {
	r io.ReadSeeker
}

func NewFetcher(r io.ReadSeeker) Fetcher {
	return &fetcher{r}
}

func (fc *fetcher) Fetch(offset int64, d Decoder) error {
	_, err := fc.r.Seek(offset, 0)
	if err != nil {
		return err
	}

	err = d.DecodeFrom(fc.r)
	if err != nil {
		return err
	}

	return nil
}
