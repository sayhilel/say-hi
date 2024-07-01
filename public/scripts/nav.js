function listenNav() {
    let index = 0;
    document.addEventListener("keydown",
        function(event) {
            // Global
            if (event.ctrlKey && event.key === "c") {
                htmx.ajax('GET', "about-me", '#terminal-window');
            }

            // Prompt
            let prompt = document.getElementById('command');
            if (prompt) {
                prompt.focus()
            }

            // Projects
            let projects = document.getElementById("projects");
            if (projects) {
                let opt = document.querySelectorAll(".selector li");

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

        });
}
