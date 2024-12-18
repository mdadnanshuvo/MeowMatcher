<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat API Project</title>
    <style>
        /* Add some basic styling for responsiveness */
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        header {
            text-align: center;
            background: #ff9800;
            color: white;
            padding: 10px 0;
        }
        nav {
            display: flex;
            justify-content: center;
            gap: 20px;
            background: #ffcc80;
            padding: 10px 0;
        }
        nav button {
            padding: 10px 20px;
            border: none;
            cursor: pointer;
            background: #fff;
            border-radius: 5px;
            font-size: 16px;
        }
        nav button.active {
            background: #ff9800;
            color: white;
        }
        section {
            margin: 20px;
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 10px;
        }
        .card {
            border: 1px solid #ddd;
            border-radius: 5px;
            overflow: hidden;
            text-align: center;
        }
        .card img {
            width: 100%;
            height: 150px;
            object-fit: cover;
        }
        .card button {
            background: #ff9800;
            color: white;
            border: none;
            padding: 5px 10px;
            margin: 5px;
            border-radius: 5px;
            cursor: pointer;
        }
    </style>
</head>
<body>
    <header>
        <h1>Cat API Project</h1>
    </header>
    <nav>
        <button class="tab" data-tab="voting">Voting</button>
        <button class="tab" data-tab="breeds">Breeds</button>
        <button class="tab" data-tab="favorites">Favorites</button>
    </nav>
    <section id="content">
        <!-- Dynamic content goes here -->
    </section>
    <script src="/static/app.js"></script>
</body>
</html>
