document.querySelector("form").addEventListener("submit", function (e) {
    const title = document.getElementById("title").value.trim();
    const content = document.getElementById("content").value.trim();
    const categories = document.querySelectorAll('input[name="category"]:checked');

    let errors = [];

    if (title === "") {
        errors.push("Title is required");
    }

    if (content === "") {
        errors.push("Content is required");
    }

    if (categories.length === 0) {
        errors.push("Select at least one category");
    }

    if (errors.length > 0) {
        e.preventDefault(); // stop form submission
        alert(errors.join("\n"));
    }
});
