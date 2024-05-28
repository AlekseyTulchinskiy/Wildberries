package main

import (
	database "consumer/pkg/database"
	handler "consumer/pkg/handler"
	memcache "consumer/pkg/memcache"
	nats_streaming_connect "consumer/pkg/nats_streaming_connect"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 1) Подключение к Postgresql
	db, err := database.DBConnect()
	if err != nil {
		log.Fatalf("Ошибка соединения с базой данных: %v", err)
	}

	// 2) Инициализация кеша cache и заполнение данными из Postgresql

	var cache = memcache.New(5*time.Minute, 15*time.Minute)
	cache.Input(db)

	// 3) Подключение к NATS Streaming серверу

	go func() {
		err := nats_streaming_connect.СonnectingNats(db, cache)
		if err != nil {
			log.Fatalf("Ошибка при подключении к NATS: %v", err)
		}
	}()
	log.Println("Consumer запущен. Ожидание сообщений...")

	// 4) Запуск сервера
	go func() {
		err = handler.ServerStart(cache)
		if err != nil {
			log.Fatalf("Ошибка при подключении к NATS: %v", err)
		}
	}()
	// Ожидание сигнала для завершения работы приложения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	log.Println("Consumer завершает работу.")

}
