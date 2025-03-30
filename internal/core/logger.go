package core

import (
    "io"
    "log"
    "os"
)


var AppLogger *log.Logger


// LoggerFactory возвращает новый логгер, который пишет в консоль
func LoggerFactory() (*log.Logger, error) {
    // Создаем Writer, который будет писать в консоль
    multiWriter := io.MultiWriter(os.Stdout)

    // Создаем новый логгер
    logger := log.New(multiWriter, "[BOT] ", log.LstdFlags|log.Lmicroseconds)

    return logger, nil
}