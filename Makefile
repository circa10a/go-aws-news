GOCMD=go
VERSION=0.1.5

# First target for travis ci
test:
	$(GOCMD) test -v

coverage:
	$(GOCMD) test -coverprofile=c.out && go tool cover -html=c.out && rm c.out
