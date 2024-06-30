function returnHome(event) {
    if (event.ctrlKey && event.key === "c") {
        htmx.ajax('GET', "about-me", '#terminal-window');
    }
}

function initializeNav() {
    let opt = document.querySelectorAll(".selector li");
    let index = 0;
    opt[index].classList.add("highlight");

    document.addEventListener("keydown", navHandle);

    function navHandle(event) {
        const projectsDiv = document.getElementById("projects");
        if (!projectsDiv) return;

        const rect = projectsDiv.getBoundingClientRect();
        const isVisible = rect.top >= 0 && rect.bottom <= (window.innerHeight || document.documentElement.clientHeight);

        if (isVisible) {
            const url = `/projects/${index}`;
            switch (event.key) {
                case "ArrowUp":
                    opt[index].classList.remove("highlight");
                    index = (index - 1 + opt.length) % opt.length;
                    opt[index].classList.add("highlight");
                    htmx.ajax('GET', url, '#project-box');
                    break;
                case "ArrowDown":
                    opt[index].classList.remove("highlight");
                    index = (index + 1) % opt.length;
                    opt[index].classList.add("highlight");
                    htmx.ajax('GET', url, '#project-box');
                    break;
            }
        }
    }
}
