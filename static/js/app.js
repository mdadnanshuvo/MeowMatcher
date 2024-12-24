document.addEventListener("DOMContentLoaded", () => {
    const content = document.getElementById("cat-content");
    const tabs = document.querySelectorAll(".tab");

    // Initialize the app with the default tab
    loadVoting();

   
// Adding Icons to Tabs Dynamically
tabs.forEach((tab) => {
    const tabType = tab.dataset.tab;

    // Add icons based on the tab type
    const iconMap = {
        voting: "/static/icons/voting.png",
        breeds: "/static/icons/breeds.png",
        favorites: "/static/icons/fav.png",
    };

    const iconSrc = iconMap[tabType];
    if (iconSrc) {
        const iconElement = document.createElement("img");
        iconElement.src = iconSrc;
        iconElement.alt = `${tabType} Icon`;
        iconElement.classList.add("tab-icon");

        // Prepend the icon to the tab
        tab.prepend(iconElement);
    }

    tab.addEventListener("click", (e) => {
        e.preventDefault();
        tabs.forEach((t) => t.classList.remove("active"));
        tab.classList.add("active");

        const selectedTab = tabType;
        if (selectedTab === "voting") loadVoting();
        else if (selectedTab === "favorites") loadFavorites();
        else if (selectedTab === "breeds") loadBreeds();
    });
});



// Adding hover placeholder functionality
let hoverPlaceholder;
function initializeHoverPlaceholder() {
    hoverPlaceholder = document.createElement("div");
    hoverPlaceholder.id = "hover-placeholder";
    hoverPlaceholder.style.position = "absolute";
    hoverPlaceholder.style.background = "rgba(0, 255, 0, 0.9)"; // Green background
    hoverPlaceholder.style.color = "white";
    hoverPlaceholder.style.padding = "5px 10px";
    hoverPlaceholder.style.borderRadius = "5px";
    hoverPlaceholder.style.fontSize = "14px";
    hoverPlaceholder.style.whiteSpace = "nowrap";
    hoverPlaceholder.style.pointerEvents = "none"; // Prevent blocking clicks
    hoverPlaceholder.style.display = "none"; // Initially hidden
    hoverPlaceholder.style.zIndex = "1000";
    hoverPlaceholder.textContent = "Click to view in gallery"; // Placeholder message
    document.body.appendChild(hoverPlaceholder);
}

function attachHoverEvents() {
    const favoritesContainer = document.getElementById("favorites-container");
    if (!favoritesContainer) return;

    // Show placeholder on mouseover
    favoritesContainer.addEventListener("mouseover", (e) => {
        if (e.target.tagName === "IMG") {
            hoverPlaceholder.style.display = "block";
        }
    });

    // Update placeholder position on mousemove
    favoritesContainer.addEventListener("mousemove", (e) => {
        if (e.target.tagName === "IMG") {
            hoverPlaceholder.style.top = `${e.pageY + 15}px`; // Offset below cursor
            hoverPlaceholder.style.left = `${e.pageX + 15}px`; // Offset to the right of cursor
        }
    });

    // Hide placeholder on mouseout
    favoritesContainer.addEventListener("mouseout", (e) => {
        if (e.target.tagName === "IMG") {
            hoverPlaceholder.style.display = "none";
        }
    });

   
}

function deleteFavorite(id) {
    let favorites = JSON.parse(localStorage.getItem("favorites")) || [];
    favorites = favorites.filter((fav) => fav.id !== id);
    localStorage.setItem("favorites", JSON.stringify(favorites));

    // Reload favorites and check for images
    loadFavorites();
    const favoritesContainer = document.getElementById("favorites-container");
    const images = favoritesContainer.querySelectorAll("img");
    if (images.length === 0) {
        hoverPlaceholder.style.display = "none"; // Hide placeholder when no images are left
    }
}



initializeHoverPlaceholder();

// Load the Voting Tab
async function loadVoting() {
    showLoader();
    try {
        const cat = await fetchRandomCat();
        if (cat) {
            renderCat(cat);
        } else {
            content.innerHTML = `
                <p>
                    <img src="/static/icons/voting.png" alt="Voting Icon" class="inline-icon">
                    Could not load cat. Try again!
                </p>`;
        }
    } catch (error) {
        content.innerHTML = `
            <p>
                <img src="/static/icons/voting.png" alt="Voting Icon" class="inline-icon">
                Failed to load the voting tab. Please try again!
            </p>`;
    }
}


    // Show Loader Animation
    function showLoader() {
        content.innerHTML = `
            <div style="text-align: center; padding: 20px;">
                <div class="spinner"></div>
                <p>Loading...</p>
            </div>
        `;
    }

    // Fetch a Random Cat
    async function fetchRandomCat() {
        try {
            const response = await fetch("/voting"); // Fetch random cat
            if (!response.ok) throw new Error("Failed to fetch cat image");
            const data = await response.json();
            return data[0];
        } catch (error) {
            console.error("Error fetching random cat:", error);
            showToast("Failed to load cat. Please try again!");
            return null;
        }
    }
    

    // Render Cat Content
    function renderCat(cat) {
        content.innerHTML = `
    <div class="card" id="cat-${cat.id}">
        <div class="image-container">
            <img src="${cat.url}" alt="Random Cat" class="cat-image">
            <div class="image-overlay">
                <!-- Heart Button -->
                <button class="heart" data-id="${cat.id}" data-cat='${JSON.stringify(cat)}'>
                    <img src="/static/icons/heart.png" alt="Heart Icon">
                </button>
                <!-- Like and Dislike Buttons -->
                <div class="right-buttons">
                    <button class="upvote" data-id="${cat.id}">
                        <img src="/static/icons/like.png" alt="Like Icon">
                    </button>
                    <button class="downvote" data-id="${cat.id}">
                        <img src="/static/icons/dislike.png" alt="Dislike Icon">
                    </button>
                </div>
            </div>
        </div>
    </div>
`;
        // Event listeners for heart, upvote, and downvote buttons
        document.querySelector(".heart").addEventListener("click", () => {
            saveFavorite(cat);
            localStorage.removeItem("votingData"); // Clear cache after action
            loadVoting();
        });
    
        document.querySelector(".upvote").addEventListener("click", () => {
            showToast("You upvoted this cat!");
        });
    
        document.querySelector(".downvote").addEventListener("click", () => {
            showToast("You downvoted this cat!");
        });
    }

    
    // Generate a unique sub_id and store it in localStorage
function getOrCreateSubID() {
    let subID = localStorage.getItem("sub_id");
    if (!subID) {
        subID = `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        localStorage.setItem("sub_id", subID);
    }
    return subID;
}


function saveFavorite(cat) {
    const url = '/add-favourites'; // Backend route for adding favorites

    // Prepare the data to send in the POST request
    const data = {
        image_id: cat.id // Send the image ID for the favorite
    };

    // Send the POST request to the backend
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data) // Send the payload as JSON
    })
    .then(response => {
        if (response.status === 409) {
            // Handle duplicate favorite (HTTP 409 Conflict)
            showToast("This image is already a favorite! ❤️");
            return response.json(); // Still parse the response for any additional data
        }
        return response.json();
    })
    .then(responseData => {
        // Handle the response from the server
        if (responseData.message) {
            showToast(responseData.message); // Show success message
        } else if (responseData.error) {
            showToast("Error: " + responseData.error); // Show error message
        }
    })
    .catch(error => {
        // Handle any network or other errors
        console.error("Error saving favorite:", error);
        showToast("An error occurred while saving the favorite. Please try again.");
    });
}




async function loadFavorites() {
    const apiUrl = '/get-favourites'; // Backend route for fetching favorites

    try {
        const response = await fetch(apiUrl, {
            method: "GET",
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (response.ok) {
            const favorites = await response.json();
            console.log("Fetched Favorites:", favorites); // Log for debugging

            if (favorites.length === 0) {
                content.innerHTML = "<p>You have no favorites yet! ❤️</p>";
                return;
            }

            // Render favorites dynamically
            content.innerHTML = `
                <div class="favorites-tab-container px-0 pb-4 mt-4 h-[380px] overflow-y-auto">
                    <div class="view-toggle flex py-4 bg-white gap-2">
                        <button class="toggle-btn active" id="grid-view" role="button" tabindex="0">
                            <img src="/static/icons/grid-view.png" alt="Grid View" class="view-icon">Grid View
                        </button>
                        <button class="toggle-btn" id="list-view" role="button" tabindex="0">
                            <img src="/static/icons/list-view.png" alt="List View" class="view-icon">List View
                        </button>
                    </div>
                    <div id="favorites-container" class="grid">
                        ${favorites
                            .map(
                                (fav, index) => `
                        <div class="grid-view-item" data-index="${index}">
                            <img src="${fav.image.url}" alt="Favorite Cat" class="object-cover w-full h-full rounded shadow clickable">
                            <div class="delete-icon-container">
                                <img src="/static/icons/delete.png" alt="Delete" class="delete-icon" data-id="${fav.id}">
                            </div>
                        </div>`
                            )
                            .join("")}
                    </div>
                    <div id="image-modal" class="modal hidden">
                        <div class="modal-content">
                            <span class="close">&times;</span>
                            <img id="modal-image" class="modal-img" src="" alt="Modal Cat">
                            <div class="modal-nav">
                                <button id="prev-image" class="nav-btn">Prev</button>
                                <button id="next-image" class="nav-btn">Next</button>
                            </div>
                        </div>
                    </div>
                </div>
            `;

            // Attach event handlers for hover, delete, and click actions
            attachHoverEvents();
            attachDeleteHandlers();
            attachClickHandlers();
            attachViewToggleHandlers();
        } else {
            const error = await response.json();
            console.error("Failed to fetch favorites:", error);
            content.innerHTML = "<p>Error fetching favorites. Please try again later.</p>";
        }
    } catch (err) {
        console.error("Network or server error:", err);
        content.innerHTML = "<p>Unable to fetch favorites. Please try again later.</p>";
    }
}




    // Load the Breeds Tab
    async function loadBreeds() {
        const content = document.getElementById("cat-content");
        content.innerHTML = `
            <div style="text-align: center; padding: 20px;">
                <div class="spinner"></div>
                <p>
                    <img src="/static/icons/breeds.png" alt="Breeds Icon" class="inline-icon">
                    Loading breeds...
                </p>
            </div>`;
        try {
            const response = await fetch("/breeds-with-images");
            if (!response.ok) throw new Error("Failed to fetch breeds data");
    
            const breeds = await response.json();
            renderBreeds(breeds);
        } catch (error) {
            content.innerHTML = `
                <p>
                    <img src="/static/icons/breeds.png" alt="Breeds Icon" class="inline-icon">
                    Failed to load breeds. Please try again!
                </p>`;
            console.error("Error loading breeds:", error);
        }
    }



    function renderBreeds(breeds) {
        const content = document.getElementById("cat-content");
    
        // Assuming only one breed is displayed
        const breed = breeds[0];
        content.innerHTML = `
            <div class="breeds-search">
                <input 
                    type="text" 
                    id="breed-search" 
                    class="search-input" 
                    placeholder="Search ${breed.name}" 
                    autocomplete="off">
                <ul id="suggestions" class="suggestions-list"></ul>
            </div>
            <div class="carousel-item">
                <div class="carousel-images-wrapper">
                    <div class="carousel-images" id="carousel-${breed.id}">
                        ${breed.images.map((image, imgIndex) => `
                            <img 
                                src="${image.url}" 
                                alt="${breed.name} Image ${imgIndex + 1}" 
                                class="carousel-img ${imgIndex === 0 ? "active" : ""}">
                        `).join("")}
                    </div>
                    <div class="carousel-nav left" id="left-nav-${breed.id}">&lt;</div>
                    <div class="carousel-nav right" id="right-nav-${breed.id}">&gt;</div>
                </div>
                <div class="carousel-dots" id="dots-${breed.id}">
                    ${breed.images.map((_, imgIndex) => `
                        <span 
                            class="dot ${imgIndex === 0 ? "active" : ""}" 
                            data-index="${imgIndex}" 
                            data-breed="${breed.id}">
                        </span>
                    `).join("")}
                </div>
                <div class="breed-details">
                    <h2>${breed.name}</h2>
                    <p><strong>Origin:</strong> ${breed.origin}</p>
                    <p><strong>Life Span:</strong> ${breed.life_span} years</p>
                    <p><strong>Temperament:</strong> ${breed.temperament}</p>
                    <p class="description">${breed.description}</p>
                    <a href="${breed.wikipedia_url}" target="_blank" class="wiki-link">Learn more on Wikipedia</a>
                </div>
            </div>
        `;
    
        initializeCarouselForBreed(breed.id, breed.images.length);
        setupSearchBar(breeds); // Initialize search bar functionality
    }
    
    
    
    function initializeCarouselForBreed(breedId, imageCount) {
        const carousel = document.getElementById(`carousel-${breedId}`);
        const dots = document.querySelectorAll(`#dots-${breedId} .dot`);
        const items = carousel.querySelectorAll(".carousel-img");
        const leftNav = document.getElementById(`left-nav-${breedId}`);
        const rightNav = document.getElementById(`right-nav-${breedId}`);
        let currentIndex = 0;
    
        const updateCarousel = (index) => {
            const translateX = -(index * 100); // Shift by 100% for each image
            carousel.style.transform = `translateX(${translateX}%)`;
            dots.forEach((dot, i) => dot.classList.toggle("active", i === index));
        };
    
        const goToNext = () => {
            currentIndex = (currentIndex + 1) % imageCount;
            updateCarousel(currentIndex);
        };
    
        const goToPrev = () => {
            currentIndex = (currentIndex - 1 + imageCount) % imageCount;
            updateCarousel(currentIndex);
        };
    
        rightNav.addEventListener("click", goToNext);
        leftNav.addEventListener("click", goToPrev);
    
        let interval = setInterval(goToNext, 3000);
    
        carousel.addEventListener("mouseenter", () => clearInterval(interval));
        carousel.addEventListener("mouseleave", () => {
            interval = setInterval(goToNext, 3000);
        });
    
        updateCarousel(currentIndex);
    }
    
   

    function setupSearchBar(breeds) {
        const searchInput = document.getElementById("breed-search");
        const suggestionsList = document.getElementById("suggestions");
    
        // Populate suggestions when the search bar is focused
        searchInput.addEventListener("focus", () => {
            renderSuggestions(breeds);
        });
    
        // Handle input typing for dynamic suggestions
        searchInput.addEventListener("input", (e) => {
            const query = e.target.value.toLowerCase();
            const filteredBreeds = breeds.filter((breed) =>
                breed.name.toLowerCase().includes(query)
            );
            renderSuggestions(filteredBreeds);
        });
    
        // Handle suggestion click
        suggestionsList.addEventListener("click", (e) => {
            if (e.target.classList.contains("suggestion-item")) {
                const selectedBreed = breeds.find((breed) => breed.name === e.target.textContent.trim());
                if (selectedBreed) renderBreeds([selectedBreed]); // Re-render with selected breed
                suggestionsList.style.display = "none"; // Hide suggestions after selection
            }
        });
    
        // Close suggestions when clicking outside
        document.addEventListener("click", (e) => {
            if (!suggestionsList.contains(e.target) && e.target !== searchInput) {
                suggestionsList.style.display = "none";
            }
        });
    
        // Helper to render suggestions
        function renderSuggestions(filteredBreeds) {
            suggestionsList.innerHTML = filteredBreeds
                .map((breed) => `<li class="suggestion-item">${breed.name}</li>`)
                .join("");
            suggestionsList.style.display = filteredBreeds.length ? "block" : "none";
        }
    }
    
   
   

// Attach Delete Handlers
function attachDeleteHandlers() {
    document.querySelectorAll(".delete-icon").forEach((icon) => {
        icon.addEventListener("click", (e) => {
            const idToDelete = e.target.dataset.id;
            deleteFavorite(idToDelete);
        });
    });
}

// Delete Favorite Functionality
function deleteFavorite(id) {
    let favorites = JSON.parse(localStorage.getItem("favorites")) || [];
    favorites = favorites.filter((fav) => fav.id !== id);
    localStorage.setItem("favorites", JSON.stringify(favorites));
    loadFavorites();
}


// Attach Click Handlers for Modal Opening
function attachClickHandlers() {
    const clickableImages = document.querySelectorAll(".clickable");
    const modal = document.getElementById("image-modal");
    const modalImage = document.getElementById("modal-image");
    const closeModal = modal.querySelector(".close");
    const prevImage = document.getElementById("prev-image");
    const nextImage = document.getElementById("next-image");

    let currentImageIndex = 0;
    const favorites = JSON.parse(localStorage.getItem("favorites")) || [];

    clickableImages.forEach((img) => {
        img.addEventListener("click", (e) => {
            currentImageIndex = parseInt(e.target.closest("[data-index]").dataset.index, 10);
            modalImage.src = favorites[currentImageIndex].url;
            modal.classList.remove("hidden");
        });
    });

    closeModal.addEventListener("click", () => modal.classList.add("hidden"));
    prevImage.addEventListener("click", () => navigateImages(-1, favorites, modalImage));
    nextImage.addEventListener("click", () => navigateImages(1, favorites, modalImage));
}

// Navigate Modal Images
function navigateImages(direction, favorites, modalImage) {
    const length = favorites.length;
    let currentImageIndex = favorites.findIndex((cat) => cat.url === modalImage.src);
    currentImageIndex = (currentImageIndex + direction + length) % length;
    modalImage.src = favorites[currentImageIndex].url;
}

function attachViewToggleHandlers() {
    const gridViewBtn = document.getElementById("grid-view");
    const listViewBtn = document.getElementById("list-view");
    const favoritesContainer = document.getElementById("favorites-container");

    // Ensure all elements are available
    if (!gridViewBtn || !listViewBtn || !favoritesContainer) {
        console.error("View toggle elements not found in the DOM.");
        return;
    }

    // Utility function to toggle views
    function toggleView(activeBtn, inactiveBtn, containerClassToAdd, containerClassToRemove) {
        activeBtn.classList.add("active");
        inactiveBtn.classList.remove("active");
        favoritesContainer.classList.add(containerClassToAdd);
        favoritesContainer.classList.remove(containerClassToRemove);
    }

    // Event listener for grid view
    gridViewBtn.addEventListener("click", () => {
        toggleView(gridViewBtn, listViewBtn, "grid", "list");
    });

    // Event listener for list view
    listViewBtn.addEventListener("click", () => {
        toggleView(listViewBtn, gridViewBtn, "list", "grid");
    });
}





    // Show Toast Notification
    function showToast(message) {
        const toast = document.createElement("div");
        toast.className = "toast";
        toast.textContent = message;
        document.body.appendChild(toast);
        setTimeout(() => toast.remove(), 2000);
    }
});
