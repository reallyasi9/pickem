package main

import (
	"context"

	"cloud.google.com/go/firestore"
)

type fscommitter struct {
	n  int
	i  int
	fc *firestore.Client
	wb *firestore.WriteBatch
}

func newFSCommitter(fc *firestore.Client, n int) *fscommitter {
	return &fscommitter{n: n, i: 0, fc: fc, wb: fc.Batch()}
}

func (c *fscommitter) Set(ctx context.Context, dr *firestore.DocumentRef, data interface{}, opts ...firestore.SetOption) error {
	c.wb.Set(dr, data, opts...)
	c.i++
	c.i %= c.n
	if c.i == 0 {
		if _, err := c.wb.Commit(ctx); err != nil {
			return err
		}
		c.wb = c.fc.Batch()
	}
	return nil
}

func (c *fscommitter) Create(ctx context.Context, dr *firestore.DocumentRef, data interface{}) error {
	c.wb.Create(dr, data)
	c.i++
	c.i %= c.n
	if c.i == 0 {
		if _, err := c.wb.Commit(ctx); err != nil {
			return err
		}
		c.wb = c.fc.Batch()
	}
	return nil
}

func (c *fscommitter) Update(ctx context.Context, dr *firestore.DocumentRef, data []firestore.Update, opts ...firestore.Precondition) error {
	c.wb.Update(dr, data, opts...)
	c.i++
	c.i %= c.n
	if c.i == 0 {
		if _, err := c.wb.Commit(ctx); err != nil {
			return err
		}
		c.wb = c.fc.Batch()
	}
	return nil
}

func (c *fscommitter) Commit(ctx context.Context) error {
	if c.i != 0 {
		if _, err := c.wb.Commit(ctx); err != nil {
			return err
		}
		c.i = 0
		c.wb = c.fc.Batch()
	}
	return nil
}
