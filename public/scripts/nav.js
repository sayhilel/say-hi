var allProjects = document.querySelectorAll(".selector li");
let currentIndex = 0;
allProjects[currentIndex].classList.add("highlight");
const totalProjects = allProjects.length;
document.addEventListener("keydown", function(event) {
    const url = `/projects/${currentIndex}`;
    switch (event.key) {
        case "ArrowUp":
            allProjects[currentIndex].classList.remove("highlight");
            currentIndex = (currentIndex - 1 + totalProjects) % totalProjects;
            allProjects[currentIndex].classList.add("highlight");
            htmx.ajax('GET', url, '#project-box');
            break;
        case "ArrowDown":
            allProjects[currentIndex].classList.remove("highlight");
            currentIndex = (currentIndex + 1) % totalProjects;
            allProjects[currentIndex].classList.add("highlight");
            htmx.ajax('GET', url, '#project-box');
            break;
    }
});
