# SendGrid Terraform Provider - Enhanced Error Handling & Documentation

## Проблема

Пользователь сталкивался с криптичными HTTP 400 ошибками при работе с teammate management:

```
Error: request failed: api response: HTTP 400: {"errors":[{"message":"invalid or unassignable scopes were given","field":"scopes"}]}
```

Эти ошибки возникали особенно часто при отмене terraform операций, вызывая каскадные сбои.

## Решение

### 1. Обновлены актуальные скоупы SendGrid (2024)

**Обновлено:**

- `sendgrid/resource_sendgrid_teammate.go` - добавлены все 200+ актуальных скоупов из API
- Использованы реальные данные из `https://api.sendgrid.com/v3/scopes`

**Добавлены новые категории скоупов:**

- `access_settings.*` - Управление настройками доступа
- `design_library.*` - Библиотека дизайнов
- `email_testing.*` - Тестирование email
- `user.webhooks.*` - Настройки webhooks
- `validations.email.*` - Валидация email
- И многие другие...

### 2. Улучшенная валидация скоупов

**Функции валидации:**

- `validateTeammateScopes()` - проверяет валидность скоупов до API вызова
- `sanitizeScopes()` - автоматически удаляет системные скоупы (`2fa_exempt`, `2fa_required`)
- Проактивная валидация на этапе `terraform plan`

**Примеры ошибок:**

```hcl
# НЕПРАВИЛЬНО - будет ошибка валидации
resource "sendgrid_teammate" "example" {
  email  = "user@example.com"
  scopes = ["mail.send", "2fa_exempt", "invalid.scope"]  # ❌
}

# ПРАВИЛЬНО
resource "sendgrid_teammate" "example" {
  email  = "user@example.com"
  scopes = ["mail.send", "templates.read"]  # ✅
}
```

### 3. Расширенная обработка ошибок

**В `sdk/errors.go`:**

- `parseErrorDetails()` - анализирует специфические типы ошибок
- `enhanceError()` - предоставляет контекстную помощь
- Специальная обработка отмены операций
- Улучшенные сообщения для разных HTTP статусов

**До:**

```
Error: request failed: api response: HTTP 400: {"errors":[...]}
```

**После:**

```
Error: Invalid or unassignable scopes provided. This can happen when:
1. Using invalid scope names (check SendGrid API documentation)
2. Your SendGrid plan doesn't support certain scopes
3. Including automatically managed scopes like '2fa_exempt' or '2fa_required'

Tip: Run 'terraform plan' first to validate your configuration.

Original error: request failed: api response: HTTP 400: {"errors":[...]}
```

### 4. Комплексная документация

**Создана документация с примерами:**

- `docs/resources/teammate.md` - полная документация с 5 сценариями использования
- `docs/troubleshooting.md` - руководство по устранению неполадок
- `examples/resources/sendgrid_teammate/` - практические примеры

**Примеры включают:**

- Базовый teammate
- Админ пользователь
- SSO пользователь
- Пользователь маркетинговой команды
- Массовое создание teammates

### 5. Тестовое покрытие

**Добавлены тесты:**

- `sendgrid/validate_scopes_test.go` - unit тесты валидации
- `sendgrid/resource_sendgrid_teammate_test.go` - acceptance тесты
- Тесты для валидных/невалидных скоупов
- Тесты для автоматических скоупов

## Файлы изменений

### Основные файлы:

1. **`sendgrid/resource_sendgrid_teammate.go`**

   - Обновлены все скоупы SendGrid (200+)
   - Добавлена валидация `validateTeammateScopes()`
   - Улучшены описания полей
   - Добавлена функция `sanitizeScopes()`

2. **`sdk/errors.go`**
   - Расширенная обработка ошибок
   - Контекстные сообщения об ошибках
   - Специальная обработка отмены операций

### Документация:

3. **`docs/resources/teammate.md`** - полная документация с примерами
4. **`docs/troubleshooting.md`** - руководство по устранению неполадок
5. **`templates/resources/teammate.md.tmpl`** - шаблон для генерации документации

### Примеры:

6. **`examples/resources/sendgrid_teammate/`** - 5 практических примеров
7. **`examples/teammate-management/main.tf`** - комплексный пример

### Тесты:

8. **`sendgrid/validate_scopes_test.go`** - unit тесты валидации
9. **`sendgrid/resource_sendgrid_teammate_test.go`** - acceptance тесты

## Ключевые улучшения

### ✅ Проактивная валидация

- Ошибки ловятся на этапе `terraform plan` вместо `terraform apply`
- Понятные сообщения об ошибках с actionable solutions

### ✅ Обработка отмены операций

- Четкие инструкции по восстановлению после прерванных операций
- Автоматическое обнаружение сценариев отмены

### ✅ Актуальные скоупы

- 200+ реальных скоупов SendGrid API
- Автоматическая фильтрация системных скоупов

### ✅ Комплексная документация

- Множественные практические примеры
- Руководство по устранению неполадок
- Best practices для предотвращения ошибок

### ✅ Тестовое покрытие

- Unit тесты для валидации
- Acceptance тесты для различных сценариев
- Проверка граничных случаев

## Результат

Теперь вместо загадочных HTTP 400 ошибок пользователи получают:

1. **Детальные объяснения** причин ошибок
2. **Конкретные шаги** для решения проблем
3. **Проактивную валидацию** для предотвращения ошибок
4. **Комплексную документацию** с примерами
5. **Улучшенную обработку** сценариев отмены

Это превращает "жесть" ситуации в понятные и решаемые проблемы с четким планом действий.
