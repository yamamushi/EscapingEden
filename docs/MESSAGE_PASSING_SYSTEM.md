# Message Passing System Documentation

## Overview

The Escaping Eden game server uses a sophisticated message passing architecture to handle communication between different system components (managers). This document explains how messages flow through the system, how errors are handled, and provides examples for novice Go programmers.

## System Architecture

The server is built around several key managers that communicate through Go channels:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   UI Windows    │    │  Connection      │    │  Game Manager   │
│   (GameWindow)  │◄──►│  Manager         │◄──►│                 │
└─────────────────┘    │  (Network Hub)   │    └─────────────────┘
                       └──────────────────┘
                              ▲    ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │  Account         │    │  Character      │
                       │  Manager         │    │  Manager        │
                       └──────────────────┘    └─────────────────┘
```

### Core Managers

1. **Connection Manager** (`network/manager.go`)
   - Central hub for all network communication
   - Routes messages between clients and other managers
   - Handles client connections and disconnections

2. **Game Manager** (`game/manager.go`)
   - Processes game logic and commands
   - Manages world state, characters, and inventory
   - Handles player actions like movement, building, digging

3. **Account Manager** (`accounts/manager.go`)
   - Manages user authentication and registration
   - Handles login/logout processes
   - Manages password resets

4. **Character Manager** (`character/manager.go`)
   - Manages character creation and data
   - Handles character name validation
   - Maintains character history and statistics

5. **UI Windows** (`ui/window/gamewindow/`)
   - Client-side interface components
   - Handle user input and display game state
   - Send commands to the game through the console

## Message Types

### Core Message Structures

All messages in the system follow a consistent pattern with these key components:

```go
type BaseMessage struct {
    Type               MessageType `json:"type"`
    SenderConsoleID    string      `json:"sender_id"`
    RecipientConsoleID string      `json:"recipient_id"`
    Data               interface{} `json:"data"`
    Error              ErrorType   `json:"error"`
}
```

### Message Categories

#### 1. Console Messages (`messages/console.go`)
Messages between the console (client interface) and the connection manager:

```go
type ConsoleMessage struct {
    Type ConsoleMessageType
    Data interface{}
}

// Examples:
Console_Message_Chat           // Chat messages
Console_Message_GameCommand    // Game commands (move, dig, build)
Console_Message_LoginResponse  // Authentication responses
```

#### 2. Window Messages (`messages/window.go`)
Messages between UI windows and the console:

```go
type WindowMessage struct {
    Type    WindowMessageType
    Command WindowMessageCommand
    Data    interface{}
}

// Examples:
WM_GameCommand        // Player actions from UI
WM_RequestLogin       // Login attempts
WM_NewPopup          // UI state changes
```

#### 3. Game Manager Messages (`messages/gameManager.go`)
Messages for game logic and world interaction:

```go
type GameManagerMessage struct {
    Type               GameManagerMessageType
    SenderConsoleID    string
    RecipientConsoleID string
    Data               interface{}
    Error              GMErrorType
}

// Examples:
GameManager_MoveCharacter     // Player movement
GameManager_BuildWallCommand  // Building actions
GameManager_RequestInventory  // Inventory requests
```

## Message Flow Examples

### Example 1: Player Building a Wall

Here's how a wall building command flows through the system:

```
1. Player Input (UI)
   ┌─────────────────┐
   │ User presses 'b'│
   │ (build command) │
   └─────────┬───────┘
             ▼
2. UI Processing (GameWindow)
   ┌─────────────────┐
   │ HandleInput()   │
   │ Creates build   │
   │ menu popup      │
   └─────────┬───────┘
             ▼
3. Material Selection
   ┌─────────────────┐
   │ User selects    │
   │ "Wood" material │
   └─────────┬───────┘
             ▼
4. Direction Selection
   ┌─────────────────┐
   │ User presses 'h'│
   │ (west direction)│
   └─────────┬───────┘
             ▼
5. Command Creation
   ┌─────────────────┐
   │ BuildWallWith   │
   │ Direction()     │
   │ creates message │
   └─────────┬───────┘
             ▼
6. Message to Console
   ┌─────────────────┐
   │ WindowMessage   │
   │ Type: WM_Game   │
   │ Command         │
   └─────────┬───────┘
             ▼
7. Console Processing
   ┌─────────────────┐
   │ Console routes  │
   │ to Connection   │
   │ Manager         │
   └─────────┬───────┘
             ▼
8. Connection Manager
   ┌─────────────────┐
   │ Routes to Game  │
   │ Manager via     │
   │ GMSendMessages  │
   └─────────┬───────┘
             ▼
9. Game Manager Processing
   ┌─────────────────┐
   │ Validates build │
   │ request, checks │
   │ inventory, etc. │
   └─────────┬───────┘
             ▼
10. Response Back
   ┌─────────────────┐
   │ Success/failure │
   │ message sent    │
   │ back to client  │
   └─────────────────┘
```

### Example 2: Error Handling Flow

When an error occurs (e.g., invalid build location):

```
Game Manager Error Processing:
┌─────────────────────────────────┐
│ 1. Validate build request       │
│    - Check if location is valid │
│    - Verify player has materials│
│    - Check building permissions │
└─────────┬───────────────────────┘
          ▼
┌─────────────────────────────────┐
│ 2. Error Detected               │
│    - Location blocked           │
│    - Create error response      │
└─────────┬───────────────────────┘
          ▼
┌─────────────────────────────────┐
│ 3. Send Error Message           │
│    Type: GM_FailedBuildWall     │
│    Data: Error details          │
└─────────┬───────────────────────┘
          ▼
┌─────────────────────────────────┐
│ 4. Connection Manager Routes    │
│    - Finds target connection    │
│    - Sends to client console    │
└─────────┬───────────────────────┘
          ▼
┌─────────────────────────────────┐
│ 5. Client UI Updates            │
│    - Shows error message        │
│    - "You can't build there."   │
└─────────────────────────────────┘
```

## Channel Communication

### How Channels Work

Go channels are used for safe communication between goroutines (concurrent processes). Each manager has dedicated channels for receiving and sending messages:

```go
// Game Manager channels
type GameManager struct {
    ReceiveChannel chan messages.GameManagerMessage       // Incoming messages
    SendChannel    chan messages.ConnectionManagerMessage // Outgoing messages
    // ... other fields
}

// Connection Manager channels  
type ConnectionManager struct {
    CMReceiveMessages chan messages.ConnectionManagerMessage // Incoming
    AMSendMessages    chan messages.AccountManagerMessage    // To Account Manager
    CMSendMessages    chan messages.CharacterManagerMessage  // To Character Manager
    GMSendMessages    chan messages.GameManagerMessage       // To Game Manager
    // ... other fields
}
```

### Channel Safety

Channels provide thread-safe communication:

```go
// Sending a message (non-blocking with goroutine)
go func() {
    message := messages.GameManagerMessage{
        Type: messages.GameManager_BuildWallCommand,
        SenderConsoleID: playerID,
        Data: buildData,
    }
    gm.SendChannel <- message
}()

// Receiving messages (blocking until message arrives)
for {
    select {
    case receivedMessage := <-gm.ReceiveChannel:
        // Process the message
        gm.handleMessage(receivedMessage)
    }
}
```

## Error Handling Patterns

### 1. Validation Errors

Before processing any command, managers validate the request:

```go
func (gm *GameManager) HandleBuildWall(msg messages.GameManagerMessage) {
    // Extract build data
    buildData := msg.Data.(messages.GameCharBuildWall)
    
    // Validate player exists
    character := gm.GetCharacter(msg.SenderConsoleID)
    if character == nil {
        gm.sendError(msg.SenderConsoleID, messages.GMError_InvalidCharacter)
        return
    }
    
    // Validate build location
    if !gm.isValidBuildLocation(buildData.DeltaX, buildData.DeltaY) {
        gm.sendBuildFailure(msg.SenderConsoleID)
        return
    }
    
    // Validate materials
    if !gm.hasRequiredMaterials(character, buildData.ItemID) {
        gm.sendError(msg.SenderConsoleID, messages.GMError_InsufficientMaterials)
        return
    }
    
    // Process successful build
    gm.processBuild(character, buildData)
}
```

### 2. Database Errors

Database operations are wrapped with error handling:

```go
func (gm *GameManager) SaveCharacter(character *Character) error {
    err := gm.DB.UpdateRecord("Characters", character)
    if err != nil {
        gm.Log.Println(logging.LogError, "Failed to save character:", err)
        
        // Send error to client
        errorMsg := messages.ConnectionManagerMessage{
            Type: messages.ConnectManager_Message_Error,
            RecipientConsoleID: character.ConsoleID,
            Data: "Failed to save character data",
        }
        gm.SendChannel <- errorMsg
        
        return err
    }
    return nil
}
```

### 3. Network Errors

Connection issues are handled gracefully:

```go
func (cm *ConnectionManager) HandleDisconnect(connection *Connection) {
    cm.mutex.Lock()
    defer cm.mutex.Unlock()
    
    // Remove from connection map
    cm.connectionMap.Delete(connection.ID)
    
    // Notify game manager of disconnect
    disconnectMsg := messages.GameManagerMessage{
        Type: messages.GameManager_NotifyDisconnect,
        Data: connection.ID,
    }
    cm.GMSendMessages <- disconnectMsg
    
    cm.Log.Println(logging.LogInfo, "Client disconnected:", connection.ID)
}
```

## Message Processing Patterns

### 1. Message Router Pattern

The Connection Manager acts as a central router:

```go
func (cm *ConnectionManager) MessageParser() {
    for {
        select {
        case managerMessage := <-cm.CMReceiveMessages:
            switch managerMessage.Type {
            case messages.ConnectManager_Message_Chat:
                cm.routeChatMessage(managerMessage)
            case messages.ConnectManager_Message_GameCommand:
                cm.routeToGameManager(managerMessage)
            case messages.ConnectManager_Message_AccountLogin:
                cm.routeToAccountManager(managerMessage)
            // ... handle other message types
            }
        }
    }
}
```

### 2. Command Processing Pattern

Game Manager processes commands through a structured approach:

```go
func (gm *GameManager) ProcessGameCommand(msg messages.GameManagerMessage) {
    switch msg.Type {
    case messages.GameManager_MoveCharacter:
        gm.handleMovement(msg)
    case messages.GameManager_BuildWallCommand:
        gm.handleBuildWall(msg)
    case messages.GameManager_DigCommand:
        gm.handleDig(msg)
    case messages.GameManager_RequestInventory:
        gm.handleInventoryRequest(msg)
    default:
        gm.Log.Println(logging.LogWarn, "Unknown command type:", msg.Type)
    }
}
```

### 3. Response Pattern

Managers send responses back through the chain:

```go
func (gm *GameManager) sendInventoryResponse(consoleID string, inventory []Item) {
    response := messages.ConnectionManagerMessage{
        Type: messages.ConnectManager_Message_GameCommandResponse,
        RecipientConsoleID: consoleID,
        Data: messages.GameMessage{
            Type: messages.GM_Inventory,
            Data: messages.GameMessageData{
                CharacterID: consoleID,
                Data: inventory,
            },
        },
    }
    gm.SendChannel <- response
}
```

## Debugging Message Flow

### 1. Enable Debug Logging

Add logging at key points to trace message flow:

```go
// In sender
gm.Log.Println(logging.LogInfo, "Sending build wall command for player:", playerID)

// In receiver  
gm.Log.Println(logging.LogInfo, "Received build wall command from:", msg.SenderConsoleID)

// In processor
gm.Log.Println(logging.LogInfo, "Processing build at location:", deltaX, deltaY)
```

### 2. Message Tracing

Use unique message IDs to trace messages through the system:

```go
type TrackedMessage struct {
    ID       string
    Type     MessageType
    Sender   string
    Data     interface{}
    Created  time.Time
}

func (gm *GameManager) sendTrackedMessage(msg TrackedMessage) {
    gm.Log.Println(logging.LogDebug, "Sending message:", msg.ID, "type:", msg.Type)
    // ... send message
}
```

### 3. Common Issues and Solutions

**Issue: Messages not reaching destination**
- Check channel buffer sizes
- Verify message routing logic
- Ensure goroutines are running

**Issue: Deadlocks**
- Use `select` statements with timeouts
- Avoid holding mutexes while sending on channels
- Use buffered channels where appropriate

**Issue: Message ordering**
- Messages on the same channel maintain order
- Use sequence numbers for critical ordering
- Consider using separate channels for different priorities

## Best Practices

### 1. Message Design

- Keep messages immutable after creation
- Use specific message types rather than generic ones
- Include all necessary context in the message
- Validate message data at boundaries

### 2. Error Handling

- Always handle channel send/receive errors
- Provide meaningful error messages to users
- Log errors with sufficient context
- Implement retry logic for transient failures

### 3. Performance

- Use buffered channels to prevent blocking
- Process messages in separate goroutines
- Batch related operations when possible
- Monitor channel queue lengths

### 4. Testing

- Mock channels for unit testing
- Test error conditions explicitly
- Verify message routing paths
- Load test with concurrent messages

## Conclusion

The message passing system in Escaping Eden provides a robust, scalable architecture for handling complex game interactions. By understanding the flow of messages between managers, developers can effectively debug issues, add new features, and maintain the system.

The key principles are:
- **Separation of Concerns**: Each manager has a specific responsibility
- **Asynchronous Communication**: Channels enable non-blocking message passing
- **Error Resilience**: Multiple layers of error handling ensure system stability
- **Traceability**: Comprehensive logging enables effective debugging

For more specific implementation details, refer to the individual manager source files and the message type definitions in the `messages/` directory.