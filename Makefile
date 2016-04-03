all: one_dns_pool

one_dns_pool: one_dns_pool.go Makefile
	go build one_dns_pool.go