package core

import (
    "fmt"
    "io"
    "log"
    "os"
)

// LoggerFactory возвращает новый логгер, который пишет в консоль и файл
func LoggerFactory(logFilePath string) (*log.Logger, error) {
    // Открываем файл для записи логов
    file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("не удалось открыть файл для записи логов: %w", err)
    }
    defer file.Close()

    // Создаем Writer, который будет писать в консоль и файл одновременно
    multiWriter := io.MultiWriter(os.Stdout, file)

    // Создаем новый логгер
    logger := log.New(multiWriter, "[BOT] ", log.LstdFlags|log.Lmicroseconds)

    return logger, nil
}