

# Documentación del Servidor "UbicaBus" | UbicaBus Server Documentation

> [!NOTE]
> There are two versions, one in English and one in Spanish, this is the English version, below is the Spanish version.

## Index

1. [Introduction](#introduction)
2. [Project Structure](#project-structure)
3. [Requirements and Dependencies](#requirements-and-dependencies)
4. [Installation and Configuration](#installation-and-configuration)
5. [Running the Server](#running-the-server)
6. [Endpoints and Features](#endpoints-and-features)
   - [HTTP with Gin](#http-with-gin)
   - [WebSockets Server](#websockets-server)
   - [MQTT Client](#mqtt-client)
7. [Architecture and Design](#architecture-and-design)
8. [Commit Formatting](#commit-formatting)
9. [Final Considerations](#final-considerations)
10. [References](#references)

---

## Introduction

The **UbicaBus** project is an example of a server developed in Go that follows a clean and modular architecture. Three key technologies have been integrated for communication:

- **HTTP**: Using the Gin framework to expose REST endpoints.
- **WebSockets**: Employing the Gorilla WebSocket library to create a WebSockets server.
- **MQTT**: Integrating Eclipse Paho to manage messages via MQTT.

This approach enables the construction of a scalable backend that is easy to maintain and adaptable to multiple communication channels.

---

## Project Structure

The code organisation is divided into layers, following a separation of concerns pattern:

```
/UbicaBus
 ├── cmd
 │    └── main.go                // Application entry point
 ├── domain
 │    └── user.go                // Example entity (domain model)
 ├── application
 │    └── user_service.go        // Business logic (application services)
 ├── infrastructure
 │    ├── delivery
 │    │    ├── http_handler.go   // HTTP endpoints with Gin
 │    │    ├── mqtt_handler.go   // MQTT client configuration and handling
 │    │    └── websocket_handler.go  // WebSockets server
 │    └── persistence
 │         └── db.go             // Database connection (MongoDB, Singleton pattern)
 └── go.mod                     // Project dependency management
```

Each layer has a specific purpose:

- **cmd:** Main file that initialises the application.
- **domain:** Definition of entities and business models.
- **application:** Contains business logic.
- **infrastructure:** Implements technical details such as delivery (HTTP, WebSockets, MQTT) and persistence (MongoDB).

---

## Requirements and Dependencies

Before running the server, ensure you have:

- **Go** (version 1.20 or later).
- **MongoDB** (running locally or accessible via URI).
- **MQTT Broker** (running locally on `tcp://localhost:1883` or adjust the configuration).

Dependencies installed via `go get`:

- Gin: `github.com/gin-gonic/gin`
- MongoDB Driver: `go.mongodb.org/mongo-driver/mongo`
- Eclipse Paho MQTT: `github.com/eclipse/paho.mqtt.golang`
- Gorilla WebSocket: `github.com/gorilla/websocket`

---

## Installation and Configuration

1. **Initialise the Go module (if not already done):**

   ```bash
   go mod init UbicaBus
   ```

2. **Download the required dependencies:**

   ```bash
   go get github.com/gin-gonic/gin
   go get go.mongodb.org/mongo-driver/mongo
   go get github.com/eclipse/paho.mqtt.golang
   go get github.com/gorilla/websocket
   ```

3. **Configure the MongoDB connection:**  
   In `infrastructure/persistence/db.go`, the connection to MongoDB is established using the Singleton pattern. Ensure that the URI `"mongodb://localhost:27017"` matches your setup.

4. **Configure the MQTT Broker:**  
   In `infrastructure/delivery/mqtt_handler.go`, the connection and subscription to the MQTT broker are configured. Ensure the broker is running and that the URL (`tcp://localhost:1883`) is correct.

---

## Running the Server

To start the server, run the following command from the project root:

```bash
go run cmd/main.go
```

The server will start:
- The **HTTP Server** on port `8080`.
- The **MQTT Client** in a goroutine, maintaining the broker connection.
- The **WebSockets Server** integrated into the HTTP endpoint.

---

## Endpoints and Features

### HTTP with Gin

The server exposes the following endpoints:

- **GET /hello**  
  Responds with a JSON message: `"Hello World"`.

- **GET /users**  
  Returns a simulated list of users using the service defined in `application/user_service.go`.

### WebSockets Server

- **GET /ws**  
  Establishes a WebSocket connection using Gorilla WebSocket.  
  - **Functionality:**  
    Once the connection is established, the server acts as an "echo server," returning any received messages.

### MQTT Client

The MQTT service runs in parallel with the HTTP server:
- Connects to the configured broker (`tcp://localhost:1883`).
- Subscribes to the topic `"my/topic"`.
- Messages received on that topic are logged to the server console.

---

## Architecture and Design

The project follows the principles of **clean architecture** and **hexagonal architecture**:

- **Separation of Concerns:**  
  Each layer (domain, application, infrastructure) is isolated, allowing modifications or extensions without affecting core business logic.

- **Singleton Pattern for Database:**  
  Used in `infrastructure/persistence/db.go` to ensure a single MongoDB connection instance.

- **Integration of Multiple Protocols:**  
  The server combines HTTP, WebSockets, and MQTT, facilitating communication across different channels.

---

## Commit Formatting

To maintain a clear and consistent change history, it is recommended to follow the **Conventional Commits** format. This format facilitates automatic generation of logs, changelogs, and semantic versioning. Some commit format examples:

- **feat:** For new features.  
  Example:  
  ```
  feat: add endpoint to retrieve users
  ```
- **fix:** For bug fixes.  
  Example:  
  ```
  fix: correct MongoDB connection error
  ```
- **docs:** For documentation updates.  
  Example:  
  ```
  docs: update endpoint documentation
  ```
- **chore:** For maintenance tasks, configurations, or dependency updates.  
  Example:  
  ```
  chore: update Gin version
  ```

**Visual Studio Code Recommendation:**  
To facilitate this format, install the **Conventional Commits** extension in Visual Studio Code. This extension helps:
- Generate commit messages following the standard.
- Validate that commits meet the proposed format.
- Improve commit consistency and quality within the development team.

You can install the extension by searching for "Conventional Commits" in the Visual Studio Code Marketplace or visiting: [Conventional Commits VSCode Extension](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits).

---

## Final Considerations

- **Scalability and Maintainability:**  
  The modularity of the project facilitates adding new endpoints, services, or integrating other protocols.

- **Testing and Development:**  
  The separation of responsibilities allows unit testing in each layer and facilitates team collaboration.

- **Production Configuration:**  
  Before deployment, review configurations (MongoDB URI, MQTT broker address, CORS handling in WebSockets, etc.) to ensure secure and efficient operation.

---

## References

- [Refactoring Guru – Clean Architecture and Design Patterns](https://refactoring.guru/)
- [Patterns.dev - Design Patterns Resources and Examples](https://www.patterns.dev/) citeturn0search0
- [Gin Web Framework](https://github.com/gin-gonic/gin) citeturn0search0
- [Gorilla WebSocket](https://github.com/gorilla/websocket) citeturn0search0
- [Eclipse Paho MQTT Client](https://github.com/eclipse/paho.mqtt.golang) citeturn0search0
- [Conventional Commits Extension for VSCode](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits) citeturn0search0


---

> [!Español]
> Esta es la version en español.
## Índice

1. [Introducción](#introducción)
2. [Estructura del Proyecto](#estructura-del-proyecto)
3. [Requisitos y Dependencias](#requisitos-y-dependencias)
4. [Instalación y Configuración](#instalación-y-configuración)
5. [Ejecutando el Servidor](#ejecutando-el-servidor)
6. [Endpoints y Funcionalidades](#endpoints-y-funcionalidades)
   - [HTTP con Gin](#http-con-gin)
   - [Servidor WebSockets](#servidor-websockets)
   - [Cliente MQTT](#cliente-mqtt)
7. [Arquitectura y Diseño](#arquitectura-y-diseño)
8. [Formato de los Commits](#formato-de-los-commits)
9. [Consideraciones Finales](#consideraciones-finales)
10. [Referencias](#referencias)

---

## Introducción

El proyecto **UbicaBus** es un ejemplo de servidor desarrollado en Go que implementa una arquitectura limpia y modular. Se han integrado tres tecnologías clave para la comunicación:

- **HTTP**: Utilizando el framework Gin para exponer endpoints REST.
- **WebSockets**: Empleando la librería Gorilla WebSocket para crear un servidor de WebSockets.
- **MQTT**: Integrando Eclipse Paho para gestionar mensajes vía MQTT.

Este enfoque permite construir un backend escalable, de fácil mantenimiento y adaptable a múltiples canales de comunicación.

---

## Estructura del Proyecto

La organización del código se ha dividido en capas, siguiendo un patrón de separación de responsabilidades:

```
/UbicaBus
 ├── cmd
 │    └── main.go                // Punto de entrada de la aplicación
 ├── domain
 │    └── user.go                // Ejemplo de entidad (modelo de dominio)
 ├── application
 │    └── user_service.go        // Lógica de negocio (servicios de aplicación)
 ├── infrastructure
 │    ├── delivery
 │    │    ├── http_handler.go   // Endpoints HTTP con Gin
 │    │    ├── mqtt_handler.go   // Configuración y manejo del cliente MQTT
 │    │    └── websocket_handler.go  // Servidor WebSockets
 │    └── persistence
 │         └── db.go             // Conexión a la base de datos (MongoDB, patrón Singleton)
 └── go.mod                     // Gestión de dependencias del proyecto
```

Cada capa cumple un propósito específico:

- **cmd:** Archivo principal que inicializa la aplicación.
- **domain:** Definición de entidades y modelos de negocio.
- **application:** Contiene la lógica de negocio.
- **infrastructure:** Implementa detalles técnicos como la entrega (HTTP, WebSockets, MQTT) y la persistencia (MongoDB).

---

## Requisitos y Dependencias

Antes de ejecutar el servidor, asegúrate de contar con:

- **Go** (versión 1.20 o superior).
- **MongoDB** (ejecutándose localmente o accesible vía URI).
- **Broker MQTT** (ejecutándose localmente en `tcp://localhost:1883` o ajusta la configuración).

Dependencias instaladas vía `go get`:

- Gin: `github.com/gin-gonic/gin`
- MongoDB Driver: `go.mongodb.org/mongo-driver/mongo`
- Eclipse Paho MQTT: `github.com/eclipse/paho.mqtt.golang`
- Gorilla WebSocket: `github.com/gorilla/websocket`

---

## Instalación y Configuración

1. **Inicializa el módulo Go (si aún no lo has hecho):**

   ```bash
   go mod init UbicaBus
   ```

2. **Descarga las dependencias necesarias:**

   ```bash
   go get github.com/gin-gonic/gin
   go get go.mongodb.org/mongo-driver/mongo
   go get github.com/eclipse/paho.mqtt.golang
   go get github.com/gorilla/websocket
   ```

3. **Configura la conexión a MongoDB:**  
   En `infrastructure/persistence/db.go` se establece la conexión a MongoDB utilizando el patrón Singleton. Verifica que la URI `"mongodb://localhost:27017"` corresponda a tu configuración.

4. **Configura el Broker MQTT:**  
   En `infrastructure/delivery/mqtt_handler.go` se configura la conexión y suscripción al broker MQTT. Asegúrate de que el broker esté corriendo y que la URL (`tcp://localhost:1883`) sea la correcta.

---

## Ejecutando el Servidor

Para iniciar el servidor, desde la raíz del proyecto ejecuta:

```bash
go run cmd/main.go
```

El servidor iniciará:
- El **HTTP Server** en el puerto `8080`.
- El **Cliente MQTT** en una goroutine, manteniendo la conexión al broker.
- El **Servidor WebSockets** integrado en el endpoint HTTP.

---

## Endpoints y Funcionalidades

### HTTP con Gin

El servidor expone los siguientes endpoints:

- **GET /hola**  
  Responde con un mensaje JSON: `"Hola Mundo"`.

- **GET /users**  
  Devuelve una lista simulada de usuarios utilizando el servicio definido en `application/user_service.go`.

### Servidor WebSockets

- **GET /ws**  
  Establece una conexión WebSocket mediante Gorilla WebSocket.  
  - **Funcionamiento:**  
    Una vez establecida la conexión, el servidor actúa como un "echo server", devolviendo cualquier mensaje recibido.

### Cliente MQTT

El servicio MQTT se ejecuta en paralelo al servidor HTTP:
- Se conecta al broker configurado (`tcp://localhost:1883`).
- Se suscribe al tópico `"mi/topico"`.
- Los mensajes recibidos en ese tópico se registran en la consola del servidor.

---

## Arquitectura y Diseño

El proyecto sigue los principios de la **arquitectura limpia** y la **arquitectura hexagonal**:

- **Separación de Responsabilidades:**  
  Cada capa (dominio, aplicación, infraestructura) está aislada, permitiendo modificar o extender funcionalidades sin afectar el núcleo del negocio.

- **Patrón Singleton para la Base de Datos:**  
  Se utiliza en `infrastructure/persistence/db.go` para garantizar una única instancia de conexión a MongoDB.

- **Integración de Múltiples Protocolos:**  
  El servidor combina HTTP, WebSockets y MQTT, facilitando la comunicación a través de diferentes canales.

---

## Formato de los Commits

Para mantener un historial de cambios claro y consistente, se recomienda seguir el formato **Conventional Commits**. Este formato facilita la generación automática de logs, changelogs y versiones semánticas. Algunos ejemplos de formatos de commit son:

- **feat:** Para nuevas funcionalidades.  
  Ejemplo:  
  ```
  feat: agregar endpoint para obtener usuarios
  ```
- **fix:** Para correcciones de errores.  
  Ejemplo:  
  ```
  fix: corregir error en conexión a MongoDB
  ```
- **docs:** Para cambios en la documentación.  
  Ejemplo:  
  ```
  docs: actualizar documentación de endpoints
  ```
- **chore:** Para tareas de mantenimiento, configuraciones o actualizaciones de dependencias.  
  Ejemplo:  
  ```
  chore: actualizar versión de Gin
  ```

**Recomendación en Visual Studio Code:**  
Para facilitar la adopción de este formato, se recomienda instalar la extensión **Conventional Commits** en Visual Studio Code. Esta extensión te ayudará a:
- Generar mensajes de commit siguiendo el estándar.
- Validar que los commits cumplan con el formato propuesto.
- Aumentar la consistencia y calidad de los commits en el equipo de desarrollo.

Puedes instalar la extensión buscando "Conventional Commits" en el Marketplace de Visual Studio Code o visitando el siguiente enlace: [Conventional Commits VSCode Extension](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits).

---

## Consideraciones Finales

- **Escalabilidad y Mantenibilidad:**  
  La modularidad del proyecto facilita la incorporación de nuevos endpoints, servicios o la integración de otros protocolos.
  
- **Pruebas y Desarrollo:**  
  La separación de responsabilidades permite realizar pruebas unitarias en cada capa y facilita la colaboración en equipo.

- **Configuración de Producción:**  
  Antes de desplegar en producción, revisa las configuraciones (URI de MongoDB, dirección del broker MQTT, manejo de CORS en WebSockets, etc.) para asegurar un funcionamiento seguro y eficiente.

---

## Referencias

- [Refactoring Guru – Arquitectura limpia y patrones de diseño](https://refactoring.guru/es) citeturn0search0
- [Patterns.dev – Recursos y ejemplos de patrones de diseño](https://www.patterns.dev/) citeturn0search0
- [Gin Web Framework](https://github.com/gin-gonic/gin) citeturn0search0
- [Gorilla WebSocket](https://github.com/gorilla/websocket) citeturn0search0
- [Eclipse Paho MQTT Client](https://github.com/eclipse/paho.mqtt.golang) citeturn0search0
- [Conventional Commits Extension for VSCode](https://marketplace.visualstudio.com/items?itemName=vivaxy.vscode-conventional-commits) citeturn0search0