document.addEventListener("DOMContentLoaded", () => {
    const content = document.getElementById("cat-content");
    const tabs = document.querySelectorAll(".tab");

    // Initialize the app with the default tab
    loadBreeds();
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

 // Show Loader Animation
 function showLoader() {
    content.innerHTML = `
        <div style="text-align: center; padding: 20px;">
            <div class="spinner"></div>
            <p>Loading...</p>
        </div>
    `;
}


initializeHoverPlaceholder();
// Functions for Voting

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
    
        // Event listener for heart button
        const heartButton = document.querySelector(".heart");
        heartButton.addEventListener("click", (event) => {
            showToast("You favorited this cat! ❤️");
            event.preventDefault();  // Prevent default behavior like opening tabs or links
            event.stopPropagation(); // Stop event bubbling, to avoid triggering parent listeners
            loadVoting();
            localStorage.removeItem("votingData"); // Clear cache after action
            saveFavorite(cat); // Save the favorite cat
        });
    
        // Event listener for upvote button
        const upvoteButton = document.querySelector(".upvote");
        upvoteButton.addEventListener("click", (event) => {
            showToast("You upvoted this cat! ❤️");
            event.preventDefault();  // Prevent default behavior if necessary
            event.stopPropagation(); // Stop event bubbling
            postVote(cat.id, 1); // 1 for upvote
            loadVoting();
        });
    
        // Event listener for downvote button
        const downvoteButton = document.querySelector(".downvote");
        downvoteButton.addEventListener("click", (event) => {
            showToast("You downvoted this cat. 😞");
            event.preventDefault();  // Prevent default behavior if necessary
            event.stopPropagation(); // Stop event bubbling
            postVote(cat.id, -1); // -1 for downvote
            loadVoting();
        });
    }
    
   
   
// Post Vote (upvote or downvote)
async function postVote(imageId, value) {
    const subId = "user_static_sub_id_12345"; // You can change this to dynamically fetch the user ID
    const apiUrl = "/vote"; // Your backend API endpoint

    const payload = {
        image_id: imageId,
        sub_id: subId,
        value: value, // 1 for upvote, -1 for downvote
    };

    try {
        const response = await fetch(apiUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });

        const data = await response.json();
        if (response.ok) {
            if (value === 1) {
             
            } else {
                
            }
        } else {
            showToast(`Error: ${data.error || "Unable to vote"}`);
        }
    } catch (error) {
        console.error("Error during voting:", error);
        showToast("An error occurred while voting. Please try again.");
    }
}


    //Functions for Breeds

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



    function renderBreeds(breeds, selectedBreed = null) {
        const content = document.getElementById("cat-content");
        const breed = selectedBreed || breeds[0]; // Use selected breed if provided, otherwise default to the first breed
    
        content.innerHTML = `
            <div class="breeds-search">
                <input 
                    type="text" 
                    id="breed-search" 
                    class="search-input" 
                    placeholder="Search breeds" 
                    autocomplete="off">
                <ul id="suggestions" class="suggestions-list"></ul>
            </div>
            <div class="carousel-item">
                <div class="carousel-images-wrapper">
                    <div class="carousel-images" id="carousel-${breed.id}">
                        ${breed.images
                            .map(
                                (image, imgIndex) => `
                            <img 
                                src="${image.url}" 
                                alt="${breed.name} Image ${imgIndex + 1}" 
                                class="carousel-img ${imgIndex === 0 ? "active" : ""}">
                        `
                            )
                            .join("")}
                    </div>
                    <div class="carousel-nav left" id="left-nav-${breed.id}">&lt;</div>
                    <div class="carousel-nav right" id="right-nav-${breed.id}">&gt;</div>
                </div>
                <div class="carousel-dots" id="dots-${breed.id}">
                    ${breed.images
                        .map(
                            (_, imgIndex) => `
                        <span 
                            class="dot ${imgIndex === 0 ? "active" : ""}" 
                            data-index="${imgIndex}" 
                            data-breed="${breed.id}">
                        </span>
                    `
                        )
                        .join("")}
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
    
        // Reinitialize the carousel for the selected breed
        initializeCarouselForBreed(breed.id, breed.images.length);
    
        // Reinitialize the search bar with the full list of breeds
        setupSearchBar(breeds);
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
        if (!searchInput.value.trim()) {
            renderSuggestions(breeds); // Show full list if input is empty
        }
        suggestionsList.style.display = "block"; // Show suggestions dropdown
    });

    // Handle input typing for dynamic suggestions
    searchInput.addEventListener("input", (e) => {
        const query = e.target.value.toLowerCase().trim();
        const filteredBreeds = query
            ? breeds.filter((breed) =>
                  breed.name.toLowerCase().includes(query)
              )
            : breeds; // Show full list if input is empty
        renderSuggestions(filteredBreeds);

        if (filteredBreeds.length) {
            suggestionsList.style.display = "block"; // Show dropdown only if there are suggestions
        } else {
            suggestionsList.style.display = "none"; // Hide dropdown if no suggestions match
        }
    });

    // Handle suggestion click
    suggestionsList.addEventListener("click", (e) => {
        if (e.target.classList.contains("suggestion-item")) {
            const selectedBreed = breeds.find(
                (breed) => breed.name === e.target.textContent.trim()
            );
            if (selectedBreed) {
                renderBreeds(breeds, selectedBreed); // Render the full list with the selected breed highlighted
                searchInput.value = ""; // Clear the search bar
                suggestionsList.style.display = "none"; // Hide suggestions
            }
        }
    });

    // Close suggestions when clicking outside
    document.addEventListener("click", (e) => {
        if (!suggestionsList.contains(e.target) && e.target !== searchInput) {
            suggestionsList.style.display = "none"; // Hide suggestions
        }
    });

    // Helper to render suggestions
    function renderSuggestions(breeds) {
        suggestionsList.innerHTML = breeds
            .map((breed) => `<li class="suggestion-item">${breed.name}</li>`)
            .join("");
    }

    // Initialize with the full list but keep suggestions hidden
    searchInput.value = ""; // Clear input on initialization
    suggestionsList.style.display = "none"; // Hide suggestions initially
}

   
   
    // Functions for fav-tav
    
    // Save a new favorite to localStorage and then to the API
    function saveFavorite(cat) {
        const url = '/add-favourites'; // Backend route for adding favorites
        const data = { image_id: cat.id }; // Send the image ID for the favorite
    
        // Now send the favorite to the backend (API)
        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => response.json())
        .then(responseData => {
            // Log the raw response data to check what you get back from the API
            console.log("API Response:", responseData);
    
            if (responseData.message) {
                // After successfully adding a favorite, get the current favorites from localStorage
                const favorites = JSON.parse(localStorage.getItem("favorites")) || [];
    
                // Manually construct the new favorite object based on the 'cat' object
                const newFavorite = {
                    id: cat.id,
                    image: {
                        url: cat.imageUrl // Assuming `cat.imageUrl` exists and contains the URL
                    }
                };
    
                // Add the new favorite to the list
                favorites.push(newFavorite);
                localStorage.setItem("favorites", JSON.stringify(favorites)); // Save to localStorage
    
                // Log the updated favorites in localStorage to confirm
                console.log("Updated favorites in localStorage:", JSON.parse(localStorage.getItem("favorites")));
    
                // Optionally, refresh the favorites UI by loading them again from localStorage
                loadFavorites();
            } else if (responseData.error) {
                showToast("Error: " + responseData.error); // Show error message
            }
        })
        .catch(error => {
            console.error("Error saving favorite:", error);
            showToast("An error occurred while saving the favorite. Please try again.");
        });
    }
    
    
    async function loadFavorites() {
        const apiUrl = '/get-favourites'; // Backend route for fetching favorites
        const content = document.getElementById("cat-content");
    
        // Try to load favorites from localStorage first
        let favorites = JSON.parse(localStorage.getItem("favorites")) || [];
    
        // If there are no valid favorites in localStorage, fetch from the backend
        if (favorites.length === 0) {
            console.log("No valid favorites in localStorage. Fetching from API...");
    
            try {
                const response = await fetch(apiUrl, {
                    method: "GET",
                    headers: {
                        "Content-Type": "application/json"
                    }
                });
    
                if (response.ok) {
                    favorites = await response.json();
                    console.log("Fetched Favorites (raw):", favorites); // Log for debugging
    
                    // Filter out favorites without valid images
                    const validFavorites = favorites.filter(fav => fav.image && fav.image.url);
                    console.log("Valid Favorites:", validFavorites); // Log valid favorites
    
                    if (validFavorites.length === 0) {
                        content.innerHTML = "<p>You have no valid favorites yet! ❤️</p>";
                        return;
                    }
    
                    // Store valid favorites in localStorage
                    localStorage.setItem("favorites", JSON.stringify(validFavorites));
    
                    // Log the favorites in localStorage after saving them
                    console.log("Favorites in localStorage after save:", JSON.parse(localStorage.getItem("favorites")));
    
                    // Update favorites with valid data
                    favorites = validFavorites;
                } else {
                    const error = await response.json();
                    console.error("Failed to fetch favorites:", error);
                    content.innerHTML = "<p>Error fetching favorites. Please try again later.</p>";
                    return;
                }
            } catch (err) {
                console.error("Network or server error:", err);
                content.innerHTML = "<p>Unable to fetch favorites. Please try again later.</p>";
                return;
            }
        }
    
        // Render valid favorites dynamically from localStorage or API
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
                    ${favorites.map((fav, index) => `
                        <div class="grid-view-item" data-id="${fav.id}" data-index="${index}">
                            <img src="${fav.image.url}" alt="Favorite Cat" class="object-cover w-full h-full rounded shadow clickable">
                            <div class="delete-icon-container">
                                <img src="/static/icons/delete.png" alt="Delete" class="delete-icon" data-id="${fav.id}">
                            </div>
                        </div>
                    `).join("")}
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
    
        // Attach event handlers for hover, delete, modal, and view toggle actions
        attachHoverEvents();
        attachDeleteHandlers();
        attachClickHandlers(favorites); // Pass the favorites to the click handlers
        attachViewToggleHandlers(); // Ensure view toggle works
    }
    
    
    
    
    
    async function deleteFavorite(favID) {
        const apiUrl = `/delete-favourites/${favID}`; // Backend route for deleting a favorite
    
        try {
            const response = await fetch(apiUrl, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                },
            });
    
            if (response.ok) {
                console.log(`Deleted Favorite ID: ${favID}`);
    
                // Remove the deleted favorite from the DOM
                const favoriteElement = document.querySelector(`.grid-view-item[data-id="${favID}"]`);
                if (favoriteElement) {
                    favoriteElement.remove();
                }
    
                // Check if there are remaining favorites
                const favoritesContainer = document.getElementById("favorites-container");
                if (!favoritesContainer || favoritesContainer.children.length === 0) {
                    const content = document.getElementById("cat-content");
                    content.innerHTML = "<p>You have no valid favorites yet! ❤️</p>";
                }
    
                showToast("Favorite deleted successfully!");
            } else {
                const error = await response.json();
                console.error("Failed to delete favorite:", error);
                showToast(`Error: ${error.error || "Failed to delete favorite. Please try again."}`);
            }
        } catch (err) {
            console.error("Error deleting favorite:", err);
            showToast("An error occurred while deleting the favorite.");
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



// Navigate Modal Images
function navigateImages(direction, favorites, modalImage) {
    const length = favorites.length;

    // Get the index of the currently displayed image
    let currentImageIndex = favorites.findIndex(fav => fav.image.url === modalImage.src);

    // Navigate to the next or previous image
    currentImageIndex = (currentImageIndex + direction + length) % length;

    // Update the modal image to the new image
    modalImage.src = favorites[currentImageIndex].image.url;
}


function attachClickHandlers() {
    const clickableImages = document.querySelectorAll(".clickable");
    const modal = document.getElementById("image-modal");
    const modalImage = document.getElementById("modal-image");
    const closeModal = modal.querySelector(".close");
    const prevImage = document.getElementById("prev-image");
    const nextImage = document.getElementById("next-image");

    // Get favorites from localStorage
    const favorites = JSON.parse(localStorage.getItem("favorites")) || [];

    if (!favorites || favorites.length === 0) {
        console.error("No favorites found to display in the modal.");
        return;
    }

    let currentImageIndex = 0;

    // Open the modal with the clicked image
    clickableImages.forEach((img, index) => {
        img.addEventListener("click", (e) => {
            currentImageIndex = parseInt(e.target.closest("[data-index]").dataset.index, 10);

            // Set the image in the modal
            modalImage.src = favorites[currentImageIndex].image.url;
            modal.classList.remove("hidden"); // Show the modal
        });
    });

    // Close the modal when clicking the 'X' button
    closeModal.addEventListener("click", () => modal.classList.add("hidden"));

    // Navigate to the previous or next image
    prevImage.addEventListener("click", () => navigateImages(-1, favorites, modalImage));
    nextImage.addEventListener("click", () => navigateImages(1, favorites, modalImage));

    // Close the modal when clicking outside the modal content area
    window.addEventListener("click", (e) => {
        if (e.target === modal) {
            modal.classList.add("hidden");
        }
    });
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


  