package dto

import "time"

type APIOptions struct {
	Verbose        bool
	Endpoint       string
	Token          string
	DocID          string
	RequestTimeout time.Duration
}

type DumpOptions struct {
	OutputFile        string
	SchemaPackageName string
	APIPackageName    string
}

type Options struct {
	APIOptions
	DumpOptions
}
