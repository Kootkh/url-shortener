package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// ------------------------------------
// Создаём объекты
// ------------------------------------

// создаём "объект конфига"
type Config struct {
	/** struct tags - метаданные, добавляемые к полям структуры.
	*	Сами по себе они ни как не влияют на поведение программы, они лишь позволяют прикрепить дополнительную информацию к полям.
	*	Но некоторые пакеты (либо вы сами) могут их считывать и учитывать в своём поведении.
	*
	*	Такие тэги добавляются после типа поля в определении структуры и представляют собой строку, заключенную в обратные ковычки.
	*
	*	Самый популярный пример - работа с json:
	*
	*	type Person struct {
	*		Name	string	`json:"name"`
	*		Email	string	`json:"email"`
	*		Age		int		`json:"age"`
	*		private	string	`json:"-"`
	*	}
	*
	*	В данном примере у каждого поля есть тег, связанный с форматом JSON. Это указывает пакету encoding/json как сериализовать
	*	или десериализовать поля.
	*
	*	Но варианты их использования ограничены лишь вашей фантазией:
	*	-	Сериализация и десериализация (маршаллинг и анмаршаллинг): выше был пример с JSON, но потенциально это может быть любой
	*		формат: yaml, bson, jsonb, xml и др.
	*	-	Валидация: некоторые библиотеки, например, go-validator, используют struct теги для определения правил валидации данных.
	*	-	Работа с БД: теги часто используются для связи полей структуры с полями таблицы в базе данных.
	*	-	Чтение конфигов: включает в себя и анмаршаллинг (как в п.1), и некоторую валидацию (как в п.2) - важно сопоставить поля
	*		структуры и поля из файла конфига, а затем проверить что все наместе, расставить дефолтные значения и т.п.
	*
	*	`
	*	yaml:"<имя соответствующего параметра в yaml файле>"
	*	env: "<название параметра при считывании из переменной окружения>"
	*	env-default: "<значение по умолчанию, если по какой-либо причине параметр будет отсутствовать в конфиге>"
	*	env-required: "true/false <обязательность параметра для запуска>"
	*	`
	**/

	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// HTTPServer описываем в виде отдельного объекта, потому что у него есть вложенные параметры
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

// ------------------------------------
// создаём функцию, которая прочитает файл с конфигом и заполнит созданные нами объекты
// ------------------------------------

// Приставка "Must" в названии функции используется (по соглашению) когда функция вместо "возврата ошибки" дожна использовать метод "паники"
// Так делать стоит в редких случаях! Инициализация конфига - один из таких случаев!

func MustLoad() *Config {
	// ------------------------------------
	// ПЕРВОЕ: откуда читаем конфиг?
	// ------------------------------------

	// в данном случае - из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")

	// Если не найдём - уроним приложение с фаталом
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Проверяем - существует ли такой файл? Если нет - падаем с фаталом
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	// Объявляем наш объект конфига
	var cfg Config

	// Считываем файл по указанному в конфиге пути
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg
}