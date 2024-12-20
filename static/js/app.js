document.addEventListener("DOMContentLoaded", () => {
    const content = document.getElementById("cat-content");
    const tabs = document.querySelectorAll(".tab");

    // Initialize the app with the default tab
    loadVoting();

    tabs.forEach((tab) => {
        tab.addEventListener("click", (e) => {
            e.preventDefault();
            tabs.forEach((t) => t.classList.remove("active"));
            tab.classList.add("active");

            const selectedTab = tab.dataset.tab;
            if (selectedTab === "voting") loadVoting();
            else if (selectedTab === "favorites") loadFavorites();
            else if (selectedTab === "breeds") loadBreeds();
        });
    });

    // Load the Voting Tab
    async function loadVoting() {
        const cachedVotingData = JSON.parse(localStorage.getItem("votingData"));
        if (cachedVotingData) {
            renderCat(cachedVotingData);
        } else {
            showLoader();
            const cat = await fetchRandomCat();
            if (cat) {
                localStorage.setItem("votingData", JSON.stringify(cat));
                renderCat(cat);
            } else {
                content.innerHTML = "<p>Could not load cat. Try again!</p>";
            }
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
            const response = await fetch("/voting");
            if (!response.ok) throw new Error("Failed to fetch cat image");
            const data = await response.json();
            return data[0];
        } catch (error) {
            showToast("Failed to load cat. Please try again!");
            return null;
        }
    }

    // Render Cat Content
    function renderCat(cat) {
        content.innerHTML = `
            <div class="card" id="cat-${cat.id}">
                <img src="${cat.url}" alt="Random Cat">
                <div class="actions">
                    <div class="left-icons">
                        <button class="heart" data-id="${cat.id}" data-cat='${JSON.stringify(cat)}'>❤️</button>
                    </div>
                    <div class="right-icons">
                        <button class="upvote" data-id="${cat.id}">👍</button>
                        <button class="downvote" data-id="${cat.id}">👎</button>
                    </div>
                </div>
            </div>
        `;

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

    // Save Cat to Favorites
    function saveFavorite(cat) {
        const favorites = JSON.parse(localStorage.getItem("favorites")) || [];
        if (!favorites.some((fav) => fav.id === cat.id)) {
            favorites.push(cat);
            localStorage.setItem("favorites", JSON.stringify(favorites));
            showToast("Added to favorites!");
        } else {
            showToast("Already in favorites!");
        }
    }

    // Load Breeds Tab
    async function loadBreeds() {
        showLoader();
        try {
            const response = await fetch("/breeds");
            if (!response.ok) throw new Error("Failed to fetch breeds data");
            const breeds = await response.json();
            renderBreeds(breeds);
        } catch (error) {
            content.innerHTML = "<p>Could not load breeds. Try again!</p>";
        }
    }

    // Render Breeds Content
    function renderBreeds(breeds) {
        content.innerHTML = `
            <div class="breeds-search">
                <input type="text" id="breed-search" placeholder="Search Breeds">
                <button id="breed-clear">X</button>
            </div>
            <div id="breeds-container" class="gallery">
                ${breeds
                    .map(
                        (breed, index) => `
                        <div class="gallery-item ${index === 0 ? 'active' : ''}" data-index="${index}">
                            <img src="${breed.image?.url || 'placeholder.jpg'}" alt="${breed.name}">
                            <div class="breed-info">
                                <h3>${breed.name}</h3>
                                <p>${breed.origin}</p>
                                <p>${breed.description}</p>
                                <a href="${breed.wikipedia_url}" target="_blank">Wikipedia</a>
                            </div>
                        </div>`
                    )
                    .join("")}
                <div class="gallery-dots">
                    ${breeds.map((_, index) => `<span class="dot ${index === 0 ? 'active' : ''}" data-index="${index}"></span>`).join("")}
                </div>
            </div>
        `;

        initializeGallery();
    }

    function initializeGallery() {
        const items = document.querySelectorAll(".gallery-item");
        const dots = document.querySelectorAll(".gallery-dots .dot");
        let currentIndex = 0;

        const updateGallery = (index) => {
            items.forEach((item, i) => item.classList.toggle("active", i === index));
            dots.forEach((dot, i) => dot.classList.toggle("active", i === index));
        };

        const autoSwipe = () => {
            currentIndex = (currentIndex + 1) % items.length;
            updateGallery(currentIndex);
        };

        let interval = setInterval(autoSwipe, 3000);

        dots.forEach((dot) => {
            dot.addEventListener("click", (e) => {
                clearInterval(interval);
                currentIndex = parseInt(e.target.dataset.index, 10);
                updateGallery(currentIndex);
                interval = setInterval(autoSwipe, 3000);
            });
        });
    }

    // Load Favorites Tab
    function loadFavorites() {
        const favorites = JSON.parse(localStorage.getItem("favorites")) || [];
        if (favorites.length === 0) {
            content.innerHTML = "<p>You have no favorites yet! ❤️</p>";
            return;
        }

        content.innerHTML = `
            <div class="favorites-tab-container px-0 pb-4 mt-4 h-[380px] overflow-y-auto">
                <div class="view-toggle flex py-4 bg-white gap-2">
                    <button class="toggle-btn active" id="grid-view" role="button" tabindex="0">
                        <img src="/static/icons/grid-view.png" alt="Grid View" class="view-icon">
                    </button>
                    <button class="toggle-btn" id="list-view" role="button" tabindex="0">
                        <img src="/static/icons/list-view.png" alt="List View" class="view-icon">
                    </button>
                </div>
                <div id="favorites-container" class="grid">
                    ${favorites
                        .map(
                            (fav, index) => `
                        <div class="grid-view-item" data-index="${index}">
                            <img src="${fav.url}" alt="Favorite Cat" class="object-cover w-full h-full rounded shadow clickable">
                            <img src="/site/icons/delete.png" alt="Delete" class="delete-icon" data-id="${fav.id}">
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

        attachDeleteHandlers();
        attachClickHandlers();
        attachViewToggleHandlers();
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

    // Attach View Toggle Handlers
    function attachViewToggleHandlers() {
        const gridViewBtn = document.getElementById("grid-view");
        const listViewBtn = document.getElementById("list-view");
        const favoritesContainer = document.getElementById("favorites-container");

        gridViewBtn.addEventListener("click", () => {
            listViewBtn.classList.remove("active");
            gridViewBtn.classList.add("active");
            favoritesContainer.classList.remove("list");
            favoritesContainer.classList.add("grid");
        });

        listViewBtn.addEventListener("click", () => {
            gridViewBtn.classList.remove("active");
            listViewBtn.classList.add("active");
            favoritesContainer.classList.remove("grid");
            favoritesContainer.classList.add("list");
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
