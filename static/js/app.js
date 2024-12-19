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
        });
    });

    // Load the Voting Tab
    async function loadVoting() {
        showLoader();
        const cat = await fetchRandomCat();
        if (cat) renderCat(cat);
        else content.innerHTML = "<p>Could not load cat. Try again!</p>";
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
            loadVoting();
        });

        document.querySelector(".upvote").addEventListener("click", () => {
            showToast("You upvoted this cat!");
            loadVoting();
        });

        document.querySelector(".downvote").addEventListener("click", () => {
            showToast("You downvoted this cat!");
            loadVoting();
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

    // Load Favorites Tab with Grid and List View Toggle
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
                    (fav) => `
                    <div class="grid-view-item">
                        <img src="${fav.url}" alt="Favorite Cat" class="object-cover w-full h-full rounded shadow">
                    </div>`
                )
                .join("")}
        </div>
    </div>
`;





        const gridViewBtn = document.getElementById("grid-view");
        const listViewBtn = document.getElementById("list-view");
        const favoritesContainer = document.getElementById("favorites-container");

        gridViewBtn.addEventListener("click", () => {
            favoritesContainer.className = "grid grid-cols-3 gap-y-1 gap-x-1";
        });

        listViewBtn.addEventListener("click", () => {
            favoritesContainer.className = "block";
            favoritesContainer.innerHTML = favorites
                .map(
                    (fav) => `
                    <div class="flex items-center gap-4 py-2">
                        <img src="${fav.url}" alt="Favorite Cat" class="w-24 h-24 object-cover">
                    </div>`
                )
                .join("");
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