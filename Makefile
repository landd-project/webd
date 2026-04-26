
PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

TARGET = webd
SRCS = $(wildcard *.go) $(wildcard */*.go) $(wildcard */*/*.go)

all: $(TARGET)

$(TARGET): $(SRCS)
	go build -o $@ .

install:
	cp $(TARGET) $(BINDIR)/

uninstall:
	rm $(BINDIR)/$(TARGET)

clean:
	rm $(TARGET)
