package recordio

import "io"

type Fetcher interface {
	Fetch(offset int64) (Record, error)
}

type fetcher struct {
	r io.ReadSeeker
}

func NewFetcher(r io.ReadSeeker) Fetcher {
	return &fetcher{r}
}

func (fc *fetcher) Fetch(offset int64) (Record, error) {
	_, err := fc.r.Seek(offset, 0)
	if err != nil {
		return Record{}, err
	}

	r := Record{}
	err = (&r).decodeFrom(fc.r)
	if err != nil {
		return Record{}, err
	}

	return r, nil
}
