---
name: whatsapp
description: "Integración con WhatsApp Cloud API de Meta"
---

## WhatsApp Cloud API

### Cómo funciona
- Usamos la API oficial de Meta for Developers
- Webhook entrante: `POST /webhook/whatsapp`
- Verificación webhook: `GET /webhook/whatsapp`
- Envío de mensajes vía API REST

### Archivos
- `backend/internal/infra/whatsapp/client.go` - Cliente HTTP para Meta API
- `backend/internal/infra/entrypoint/rest/whatsapp_controller.go` - Webhook endpoints
- `backend/internal/application/whatsapp_service.go` - Lógica de negocio
- `backend/internal/domain/whatsapp_conversation.go` - Modelos

### Tablas DB
- `whatsapp_conversations` - Conversación activa por cliente
- `whatsapp_messages` - Historial de mensajes
- `chatbot_sessions` - Sesiones del AI agent

### Flujo de mensaje entrante
1. Meta envía POST a `/webhook/whatsapp`
2. Controlador valida token de verificación
3. Guarda mensaje en `whatsapp_messages`
4. Si la conversación está en modo AI → envía a OpenAI
5. Si no, queda en bandeja para humano
6. Respuesta (AI o humana) se envía vía Meta API

### Envío de notificaciones
- Cuando una orden cambia a "ready", se envía WhatsApp
- Cuando una orden se entrega, se envía WhatsApp
- Usar templates aprobados por Meta para mensajes proactivos

### Configuración necesaria
- `WHATSAPP_PHONE_NUMBER_ID` - ID del número en Meta
- `WHATSAPP_ACCESS_TOKEN` - Token de acceso
- `WHATSAPP_WEBHOOK_SECRET` - Token de verificación
- `WHATSAPP_API_VERSION` - v21.0 (actual)
