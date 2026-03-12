# Plan de Implementación: Guardar contenido en DB

## Contexto

El resource-service actualmente guarda el contenido scrapeado en archivos `.md` en el filesystem. El objetivo es migrar este almacenamiento a la base de datos PostgreSQL.

### Estado actual
- **resource-service**: Scraping → guarda `.md` en `files/` → responde por RabbitMQ
- **Backend-service**: Consume respuesta → actualiza `resources` table
- **Tabla `resources`**: id, url, title, excerpt, language, status_id, timestamps

### Requisitos
- Tamaño de artículos: 50KB-500KB (mediano)
- Eliminar almacenamiento en filesystem (solo DB)

---

## Plan de Implementación

### FASE 1: Base de datos

**1.1 Nueva migración** (`backend-service/migrations/00008_add_resource_contents.sql`):

```sql
CREATE TABLE IF NOT EXISTS resource_contents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  resource_id UUID UNIQUE REFERENCES resources(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_resource_contents_resource_id ON resource_contents(resource_id);
```

---

### FASE 2: Backend-service

**2.1 Nuevo modelo** (`internal/resource_content/content.go`):

```go
package resource_content

import (
	"time"

	"github.com/google/uuid"
)

type ResourceContent struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	ResourceID  uuid.UUID `gorm:"uniqueIndex;not null"`
	Content     string    `gorm:"type:text;not null"`
	CreatedAt   time.Time
}
```

**2.2 Actualizar `ResponseMessage`** (`internal/messaging/response_consumer.go`):

- Agregar campo `Content string` al struct `ResponseMessage`

**2.3 Nuevo método en service** (`internal/resource/service.go`):

- `SaveContent(resourceID uuid.UUID, content string) error`

**2.4 Actualizar handler de respuesta** (`internal/messaging/response_consumer.go`):

- Llamar al nuevo método para guardar el contenido después de actualizar la metadata

---

### FASE 3: Resource-service

**3.1 Actualizar `ResponseMessage`** (`internal/messaging/publisher.go`):

- Agregar campo `Content string` al struct

**3.2 Modificar scraper** (`internal/scraping/scraper.go`):

- Eliminar escritura a archivo
- Retornar el contenido como string (usar `strings.Builder`)

**3.3 Actualizar flujo de scraping**:

- El método `Scrape()` debe retornar el contenido en memoria en lugar de guardarlo en disco

---

## Cambios en el flujo

| Antes | Después |
|-------|---------|
| Scraping → archivo `.md` | Scraping → contenido en memoria |
| Mensaje: {title, excerpt, language} | Mensaje: {title, excerpt, language, **content**} |
| Backend actualiza metadata | Backend actualiza metadata **+ guarda content** |

---

## Consideraciones técnicas

- **PostgreSQL TEXT**: acepta hasta 1GB, los 500KB son perfectamente manejables
- **Mensajes RabbitMQ**: 500KB en JSON es aceptable
- **Eliminación de archivos**: Se elimina complejidad de filesystem y se centraliza todo en DB

---

## Orden de implementación recomendado

1. Migración de base de datos (FASE 1)
2. Backend-service: modelo + método service + actualizar handler (FASE 2)
3. Resource-service: actualizar mensaje + modificar scraper (FASE 3)
