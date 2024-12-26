# catAPIProject

## Project Description

**MeowMatcher** is a dynamic and fully responsive web platform built using **The Cat API**. The platform is designed to provide an interactive and enjoyable experience for cat lovers by integrating a wide range of cat-related data, including random cat images, detailed breed information, and fun cat facts. It is tailored for developers and cat enthusiasts who want to incorporate real-time cat data into their applications, websites, or mobile platforms. 

With a rich collection of over 60,000 images, diverse cat breeds, and interactive functionalities, **MeowMatcher** delivers a seamless experience across devices. The project allows users to explore a variety of cat images, discover different breeds, and interact with the content through features such as voting and favoriting.

The **MeowMatcher** platform is focused on:

- **Providing a broad selection of cat images**: With over 60,000 images available through **The Cat API**, users can enjoy randomly generated cat images or explore breed-specific images in a responsive and engaging interface.
- **Showcasing detailed cat breed profiles**: Get in-depth information about various cat breeds, including physical characteristics, temperament, lifespan, and origin, making it easy for users to learn about different types of cats.
- **Enabling user interaction**: The platform allows users to vote on and favorite cat images, with features such as modals and interactive tabs to enhance user engagement. 
- **Supporting multiple access levels**: Whether you're using the free tier or a commercial plan, **MeowMatcher** offers both free and paid API access, with scalable functionality to suit the needs of both casual users and businesses.

The platform is designed with advanced front-end features to ensure a fluid experience, with a responsive layout, modals for managing favorites, and other interactive UI elements to make exploring cat data both fun and intuitive.

## Key Features

- **API Integration with The Cat API**: Effortlessly fetch random or breed-specific cat images, detailed breed profiles, and cat facts from **The Cat API**.
- **Responsive Web Design**: Enjoy a smooth, mobile-friendly experience across all devices. The site adapts seamlessly to different screen sizes and orientations.
- **Interactive Favorite Tab**: Users can favorite images and interact with them in a dedicated tab, enhanced by modal pop-ups that provide more detailed information about each image.
- **Voting System**: Users can vote on their favorite cat images, contributing to the community's top-rated cat images.
- **Advanced Modal Features**: The favorite tab includes modals that allow users to view detailed image information, such as breed data and voting statistics.
- **Real-Time Data Updates**: Retrieve updated images and breed information in real time, ensuring that the content remains fresh and engaging for all users.
- **Seamless Experience Across Devices**: Fully optimized for desktop, tablet, and mobile users, ensuring the interface remains intuitive and easy to navigate no matter the platform.
- **Free and Pro Access Options**: The project supports both a free API tier with limited requests and premium features for businesses using the data commercially. 

---

This version reflects the advanced features like modals, the responsive design, and other interactive functionalities. Let me know if you need more specific changes or further adjustments to the descriptions!


## Project Structure

📁 **MeowMatcher/**  
├── 📁 **caches/**  
│   ├── 📄 `cache.go`  
│   ├── 📄 `cache_test.go`  
│    
├── 📁 **conf/**  
│   └── 📄 `app.conf`  
├── 📁 **controllers/**  
│   ├── 📄 `cat_controller_test.go`  
│   └── 📄 `cat.go`  
├── 📁 **routers/**  
│   ├── 📄 `router_test.go`  
│   └── 📄 `router.go`  
├── 📁 **static/**  
│   ├── 📁 **css/**  
│   │   └── 📄 `style.css`  
│   ├── 📁 **icons/**  
│   └── 📁 **js/**  
├── 📁 **views/**  
├── 📄 `coverage.out`  
├── 📄 `go.mod`  
├── 📄 `go.sum`  
├── 📄 `main.go`  
├── 📄 `main_test.go`  
├── 📄 `MeowMatcher.exe`  
└── 📄 `README.md`  

---

## How to Clone and Run the Project

### 1. **Clone the Repository**

To get started, you'll need to clone the **MeowMatcher** repository to your local machine. Open a terminal and run the following command:

```bash
git clone https://github.com/mdadnanshuvo/MeowMatcher.git
```
Once the repository has been cloned, navigate to the project directory:

```bash
cd MeowMatcher
```
### 2. **Install Go**

Make sure you have Go installed on your system. You can check your Go version by running the following command:

```bash
go version
```

If Go is not installed, download and install it from the official Go website. Follow the installation instructions based on your operating system.

### 3. **Install Beego**
The MeowMatcher project uses the Beego web framework. To install Beego, run the following command:

```bash
go get github.com/astaxie/beego
```
Additionally, you'll need Bee (Beego's command-line tool) for running tasks like creating and managing Beego applications. Install Bee by running:
```bash
go get github.com/beego/bee
```

### 4. **Install Project Dependencies**

After installing Go and Beego, you need to install all the necessary dependencies for the project. Run the following command to retrieve and install all dependencies specified in the go.mod file:

```bash
go mod tidy
```
This will ensure that all required libraries and modules are downloaded and set up properly.

### 5. **Configure the Application**

The **MeowMatcher** application uses **The Cat API** to fetch images, breed information, and cat facts. You must configure your **API Key** to access the data:

1. Sign up for an API key at [The Cat API](https://thecatapi.com/).
2. After obtaining the API key, open the `conf/app.conf` file.
3. Add your **API Key** in the configuration file like this:

```conf
api_key=YOUR_API_KEY_HERE
```

Make sure to replace YOUR_API_KEY_HERE with the actual API key you received from The Cat API.


### 6. **Run the Application Locally**

To start the **MeowMatcher** application locally, use the following command:

```bash
bee run
```

This command will start the Beego web server. Once the server is running, you can access the application in your browser at:

```bash
http://localhost:8080
```


### Explanation:
- The command `bee run` is wrapped in a **bash code block** for terminal commands.
- The URL `http://localhost:8080` is wrapped in an **arduino code block**, though in this context, it's just for styling the URL. You could also use a simple code block (`code` or `text`) for the URL, but I've kept it as you initially indicated.
