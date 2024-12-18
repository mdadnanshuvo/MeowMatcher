document.addEventListener("DOMContentLoaded", () => {
    const content = document.getElementById("cat-content");
    const tabs = document.querySelectorAll(".tab");

    // Default Tab
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
                        <button class="heart" data-id="${cat.id}">❤️</button>
                    </div>
                    <div class="right-icons">
                        <button class="upvote" data-id="${cat.id}">👍</button>
                        <button class="downvote" data-id="${cat.id}">👎</button>
                    </div>
                </div>
            </div>
        `;
    
        // Button Event Listeners
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

    // Load Favorites Tab
    function loadFavorites() {
        const favorites = JSON.parse(localStorage.getItem("favorites")) || [];
        if (favorites.length === 0) {
            content.innerHTML = "<p>You have no favs yet! ❤️</p>";
        } else {
            content.innerHTML = favorites.map((cat) => `
                <div class="card">
                    <img src="${cat.url}" alt="Favorite Cat">
                </div>
            `).join("");
        }
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
