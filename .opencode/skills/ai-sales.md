---
name: ai-sales
description: "Agente de ventas IA híbrido para atención al cliente"
---

## AI Sales Agent

### Comportamiento
- **Híbrido**: AI atende primeras interacciones, escala a humano cuando es necesario
- Usa OpenAI GPT-4o-mini
- Opera vía WhatsApp (principal) y web chat

### Archivos
- `backend/internal/application/ai_agent_service.go` - Lógica del agente
- `backend/internal/infra/ai/openai_client.go` - Cliente OpenAI
- `backend/internal/domain/chatbot_session.go` - Modelo de sesión
- `backend/internal/infra/repository/database/chatbot_repository.go` - Persistencia

### Cuándo escala a humano
- Cliente pide descuento
- Cliente está enojado/reclamo
- Cliente pide diseño personalizado complejo
- Cliente pide información que el AI no sabe
- Cliente pide hablar con un humano explícitamente

### Cuándo cierra la venta el AI
- Cliente elige un producto del catálogo
- Cliente acepta cotización
- Pedido simple (cantidad, material, precio estándar)

### Prompt del agente (resumen)
Eres un vendedor de stickers por WhatsApp. Tu objetivo es:
1. Saludar y preguntar qué necesita
2. Recomendar productos según su tipo de negocio
3. Cotizar precios y cantidades
4. Cerrar la venta si es posible
5. Escalar a humano si es necesario

Sé amable, directo, y útil. No inventes precios ni productos.
