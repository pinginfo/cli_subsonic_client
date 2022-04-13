# The compiler
CC = g++

# Compiler flags
CFLAGS = -std=c++20

# Target
TARGET = main

all: run

run: build
	./$(TARGET)

build: $(TARGET)

$(TARGET): $(TARGET)
	$(CC) $(CFLAGS) -o $(TARGET) src/$(TARGET).cpp

clean:
	$(RM) $(TARGET)