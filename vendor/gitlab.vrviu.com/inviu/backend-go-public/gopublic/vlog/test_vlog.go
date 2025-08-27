package vlog

import (
	"fmt"
	"log"
	"os"
)

// devlopment config log test
func DevelopTest() {
	config := NewDevelopmentConfig("test_vlog", "/tmp/devlop.log")
	logger := NewLogger(config)

	logger.Print("develop logger print")
	logger.Debug("develop logger debug")
	logger.Info("develop logger info")
	logger.Warn("develop logger warn")
	logger.Error("develop logger error")
	//logger.Fatal("develop logger fatal")
	//logger.Panic("develop logger panic")

	logger.With("test_field_string", "test_value_string").Debug("develop logger debug")
	logger.With("test_field_int", 100).Debug("develop logger debug")
}

// production config log test
func ProductTest() {
	config := NewProductionConfig("test_vlog", "/tmp/product.log")
	logger := NewLogger(config)

	logger.Debug("product logger print")
	logger.Debug("product logger debug")
	logger.Info("product logger info")
	logger.Warn("product logger warn")
	logger.Error("product logger error")
	//logger.Fatal("product logger fatal")
	//logger.Panic("product logger panic")

	logger.With("test_field_string", "test_value_string").Debug("product logger debug")
	logger.With("test_field_int", 100).Debug("product logger debug")
}

// config file log test
func ConfigTest() {
	gopath := os.Getenv("GOPATH")
	filename := gopath + "/src/public/vlog/logger.json"
	config, err := NewConfig(filename)
	if err != nil {
		fmt.Println("init logger config file failed, file:", filename, "err:", err)
		return
	}
	logger := NewLogger(config)

	logger.Debug("config logger print")
	logger.Debug("config logger debug")
	logger.Info("config logger info")
	logger.Warn("config logger warn")
	logger.Error("config logger error")
	//logger.Fatal("config logger fatal")
	//logger.Panic("config logger panic")

	logger.With("test_field_string", "test_value_string").Debug("config logger debug")
	logger.With("test_field_int", 100).Debug("config logger debug")
}

// std log test
func StdTest() {
	config := NewDevelopmentConfig("test_vlog", "/tmp/std.log")
	logger := NewLogger(config)

	RedirectStdLog(logger)
	log.Print("std logger Print")
	//log.Fatal("test")
	//log.Panic("test")
}

// global log test
func LogTest() {
	gopath := os.Getenv("GOPATH")
	filename := gopath + "/src/public/vlog/logger.json"
	config, err := NewConfig(filename)
	if err != nil {
		fmt.Println("init logger config file failed, file:", filename, "err:", err)
		return
	}
	SetConfig(config)

	Debug("log logger print")
	Debug("log logger debug")
	Info("log logger info")
	Warn("log logger warn")
	Error("log logger error")
	//Fatal("log logger fatal")
	//Panic("log logger panic")

	With("test_field_string", "test_value_string").Debug("log logger debug")
	With("test_field_int", 100).Debug("log logger debug")
}
