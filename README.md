# Groupie-Tracker 🎸

[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR-USERNAME/Groupie-Tracker)](https://goreportcard.com/report/github.com/YOUR-USERNAME/Groupie-Tracker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A music artist tracking web application built with Go that displays band information, concert dates, and locations using data from the [Groupie Trackers API](https://groupietrackers.herokuapp.com/api).

![Screenshot](https://raw.githubusercontent.com/BjornMikael/images/main/Groupie-Tracker+Screenshot.png?text=Groupie-Tracker+Screenshot)

---

## Features ✨

- 🎨 **Responsive artist card grid**
- 🔍 **Live search with suggestions**
- 📅 **Concert date tracking**
- 📍 **Location information**
- 🎤 **Detailed artist profiles**
- 🛡️ **Custom error pages (404, 500)**
- ⚡ **Auto-refreshing data**
- 📱 **Mobile-friendly design**

---

## Technologies Used 🛠️

- **Backend**: Go (Golang)
- **Frontend**: HTML5, CSS3, JavaScript
- **API**: [Groupie Trackers API](https://groupietrackers.herokuapp.com/api)
- **Styling**: CSS Grid/Flexbox, Google Fonts

---

## Installation 📦

### Prerequisites
- Go 1.20+ → [Install Guide](https://go.dev/doc/install)

### Running Locally
```bash
# Clone repository
git clone https://github.com/YOUR-USERNAME/Groupie-Tracker.git
cd Groupie-Tracker

# Install dependencies
go mod download

# Start server
go run main.go
```
Visit **[http://localhost:8080](http://localhost:8080)** in your browser 🚀

---

## Usage 🖥️

### **Home Page**
✅ Browse artists in a responsive grid  
✅ Search using the live search bar  
✅ Click artist cards for detailed views  

### **Artist Page**
✅ View members and creation date  
✅ See all concert locations  
✅ Check upcoming/past dates  
✅ Navigate back with Home button  

### **Error Pages**
✅ Custom-designed **404** and **500** pages  
✅ Easy navigation back to safety  

---

## API Endpoints 📡

| Endpoint        | Description                          |
|---------------|----------------------------------|
| `/`          | Home page with artist grid       |
| `/artist/{id}` | Individual artist details       |
| `/search`    | JSON search suggestions         |

---

## Configuration ⚙️

Modify `main.go` for:
- **Port configuration** (default: `8080`)
- **Data refresh interval** (default: `1 minute`)
- **Custom styling** in `static/style.css`

---

## Contributing 🤝

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create your feature branch**:  
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Commit your changes**:  
   ```bash
   git commit -m 'Add amazing feature'
   ```
4. **Push to the branch**:  
   ```bash
   git push origin feature/amazing-feature
   ```
5. **Open a Pull Request**

---

## License 📄

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## Acknowledgments 🙏

- **Seán David Middleton** → [GitHub](https://github.com/middsea) for great companionship, friendship, support and co-working
- Data provided by **Groupie Trackers API**
- Inspired by **music tracking applications**
- **Google Fonts** for typography
- **GitHub Community** for support
