# catAPIProject

## Project Description

**catAPIProject** is a robust and scalable platform built on **The Cat API**, offering a wide range of cat-related data to enhance user experience in applications. The project utilizes data from **The Cat API** to deliver random cat images, detailed breed information, and fun facts about cats. With over 60,000 images, breed data, and real-time API interactions, this project allows developers to integrate cat-related content into their services.

The key focus of the **catAPIProject** is to:

- **Provide access to a vast collection of cat images**: With over 60,000 images available, you can easily access random or breed-specific cat images.
- **Retrieve detailed cat breed information**: Detailed data about various breeds, including characteristics, origin, temperament, and lifespan.
- **Integrate seamlessly with other systems**: By offering API-based interactions, it allows you to build scalable and flexible applications using real-time cat data.
- **Support voting and favoriting**: Users can vote on and favorite their preferred images.
- **Support for both free and paid access**: The API offers a free tier with limited requests as well as premium features for commercial use.

This project is designed for developers and cat enthusiasts who want to integrate cat data into their applications, websites, or mobile platforms.

## Key Features

- **API Integration**: Seamless integration with **The Cat API** to get images, facts, and breed information.
- **Real-Time Data**: Fetch random or breed-specific cat images using the provided API endpoints.
- **Extensive Breed Information**: Retrieve detailed breed descriptions, characteristics, and origin information.
- **Voting and Favoriting**: Allow users to vote for or favorite images and manage those actions.
- **Flexible Pricing Options**: Access both free and pro-tier features based on project requirements.
- **Scale Efficiently**: Use real-time webhooks and high-resolution images for businesses and enterprise solutions.

---

## Pricing Plans

**catAPIProject** offers various pricing tiers for different needs:

- **Free Plan** ($0.00/month):
  - 10,000 requests per month
  - Access to images, breed information, and cat facts
  - Commercial license included
  - Useful for learning and small projects

- **Pro Plan** ($10.00/month):
  - 100,000 requests per month
  - Includes real-time web-hooks (coming soon)
  - Access to additional features like medical data, detailed breed info, and more

- **Enterprise Plan** (Price on Request):
  - Unlimited requests
  - Premium images and high-resolution content
  - Video streams and trend data
  - Custom solutions for your business needs

---

## Project Structure

📁 **catAPIProject/**  
├── 📁 **caches/**  
│   ├── 📄 `cache.go`  
│   ├── 📄 `cache_test.go`  
│   ├── 📄 `channel.go`  
│   └── 📄 `channel_test.go`  
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

To get started, you'll need to clone the repository to your local machine. Open a terminal and run the following command:

```bash
git clone https://github.com/yourusername/catAPIProject.git
