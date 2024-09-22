package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/sagernet/sing-box/option"
)

func GenerateSRSList(ruleSet RuleSet, fileName string) error {
	// Переводим итоговый rule-set в json
	jsonData, err := json.Marshal(ruleSet)
	if err != nil {
		fmt.Println("Ошибка маршализации в JSON:", err)
	}

	// Создаём переменную S-B для хранения rule-set'ов
	var plainRuleSetCompat option.PlainRuleSetCompat

	// Конвертируем полученный json функцией sing-box'а
	if plainRuleSetCompat.UnmarshalJSON(jsonData) != nil {
		return fmt.Errorf("json ruleset unmarshalization error: %v", err)
	}
	// Проверяем версию rule-set
	plainRuleSetCompat.Upgrade()

	if _, err := os.Stat("./rules"); os.IsNotExist(err) {
		// Если её нет, создаем
		err := os.MkdirAll("./rules", os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory '%s': %v", "./rules", err)
		}
		fmt.Printf("the directory '%s' was missing, but it was created:", "./rules")
	}

	// Создаём .srs файл
	RuleSetSrs, err := os.Create("./rules/" + fileName + ".srs")
	if err != nil {
		return fmt.Errorf("cannot create .srs file: %v", err)
	}
	defer RuleSetSrs.Close()

	// Пишем в .srs файл
	if err := srs.Write(RuleSetSrs, plainRuleSetCompat.Options); err != nil {
		return fmt.Errorf("cannot write into .srs file: %v", err)
	}

	return nil
}
