# 🍫 Fudge Configs

## Fudge having to manage the same config everywhere! 😤

Are you tired of copying the same configuration to 47 different files across your projects? Do you find yourself muttering "oh fudge it" every time you need to update a config that somehow needs to exist in 12 different places?

**Welcome to your salvation!** 🎉

This delightfully frustrating TUI app built with Go and BubbleTea will help you stop pulling your hair out over configuration management. Because let's be honest - life's too short to manually sync the same config across a gazillion files.

## 🤖 The MCP Madness That Started It All

This tool was born from the very real pain of managing **MCP (Model Context Protocol) configurations** across multiple AI tools. You know the drill:

- 🤖 **Claude Desktop** needs your MCP servers in `~/.config/claude/claude_desktop_config.json`
- 🤖 **Cursor** wants them in its own special place
- 🤖 **Your custom AI setup** has yet another config file
- 🤖 **That new AI tool you're trying** also needs the same servers configured

And every time you:

- ✨ Add a new MCP server (like a file browser, or calculator, or web search)
- 🔧 Update a server configuration
- 🗑️ Remove a server you're not using anymore

You have to manually update ALL of these files. Miss one? Enjoy debugging why your AI tool can't access that server! 🎭

**This madness ends here.** Configure your MCP servers once, sync everywhere, and get back to building cool stuff with AI instead of wrestling with configs.

## 🎯 What This Little Miracle Does

- **🔧 Open config**: Cracks open your master config file in whatever editor you've got (`$EDITOR`) so you can actually make changes like a civilized human
- **➕ Add a Sync file**: Tell this thing which config files need to stay in sync (because apparently they can't just behave themselves)
- **🔄 Sync config**: The magic happens here - your master config gets applied to ALL your sync files. No more copy-paste madness!

## 🚀 Getting This Party Started

```bash
go build -o fudge-configs .
```

_That's it. Really. We're not monsters here._

## 🎮 How to Use This Beauty

Fire it up:

```bash
./fudge-configs
```

### 🕹️ Navigation (It's Not Rocket Science)

- **Arrow keys** or **j/k** to navigate (vim users, we see you)
- **Enter** or **Space** to select something (shocking, we know)
- **Ctrl+C** or **q** to escape this madness
- **Esc** to go back when you inevitably get confused

### 📁 Where Your Configs Live

Your precious files hang out here:

- **macOS/Linux**: `$XDG_CONFIG_HOME/fudge_configs/` (or `~/.config/fudge_configs/` if you're basic)

#### config.json - The Boss File 👑

This is your master config that rules them all. If it doesn't exist, we'll create one for you because we're nice like that:

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-filesystem",
        "/path/to/allowed/files"
      ]
    },
    "brave-search": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-brave-search"],
      "env": {
        "BRAVE_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

_Customize this structure for whatever configs you're managing - MCP servers, database connections, API endpoints, you name it!_

#### sync.json - The Hit List 📝

This file knows which other files need to get in line:

```json
{
  "paths": [
    "~/.config/claude/claude_desktop_config.json",
    "~/.config/cursor/config.json",
    "/path/to/your/custom/ai/config.json"
  ]
}
```

### 🪄 The Sync Magic Explained

When you hit "Sync config" (the good stuff):

1. We read your boss file (`config.json`)
2. For each poor soul in your sync list:
   - We read their current sad state
   - We merge in the boss config (because boss always wins)
   - We keep any unique stuff they had (we're not completely heartless)
   - We write it back and move on to the next victim

### 🎪 Example Workflow (The Happy Path)

1. **🎭 Set up your empire**: Hit "Open config" and craft your masterpiece in `config.json`
2. **📋 Recruit your minions**: Use "Add a Sync file" to tell us which files need to fall in line
3. **⚡ Rule with an iron fist**: Hit "Sync config" and watch the magic happen

_No more "wait, did I update that MCP server in the other 6 AI tool configs?" anxiety attacks!_

### 🤖 Real-World MCP Example

Let's say you want to add a new calculator MCP server to all your AI tools:

1. **🔧 Open your master config** and add:

```json
{
  "mcpServers": {
    "calculator": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-calculator"]
    }
  }
}
```

2. **⚡ Hit sync** and boom! Calculator access is now available in:
   - Claude Desktop ✅
   - Cursor ✅
   - Your custom AI setup ✅
   - That new AI tool you're trying ✅

_One change, infinite happiness._ 🎉

## 🛠️ What You Need to Not Hate Your Life

- Go 1.21 or later (because we're not living in the stone age)
- A working `$EDITOR` environment variable (we'll use `vi` if you forgot, but come on...)

## 🎨 Built With Love (and Mild Frustration)

- [BubbleTea](https://github.com/charmbracelet/bubbletea) - Because terminal UIs should be delightful
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Making things pretty since forever

---

_Life's too short for config management hell. Go forth and sync responsibly!_ ✨
