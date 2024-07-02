function listenNav() {
    let index = 0;
    document.addEventListener("keydown",
        function(event) {
            // Global
            if (event.ctrlKey && event.key === "c") {
                htmx.ajax('GET', "landing", '#terminal-window');
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

// Thanks w3schools for this
function dragElement(elmnt) {
    var pos1 = 0, pos2 = 0, pos3 = 0, pos4 = 0;
    document.getElementById(elmnt.id + "-decorations").onmousedown = dragMouseDown;

    function dragMouseDown(e) {
        e = e || window;
        e.preventDefault();
        // get the mouse cursor position at startup:
        pos3 = e.clientX;
        pos4 = e.clientY;
        document.onmouseup = closeDragElement;
        // call a function whenever the cursor moves:
        document.onmousemove = elementDrag;
    }

    function elementDrag(e) {
        e = e || window;
        e.preventDefault();
        // calculate the new cursor position:
        pos1 = pos3 - e.clientX;
        pos2 = pos4 - e.clientY;
        pos3 = e.clientX;
        pos4 = e.clientY;
        // set the element's new position:
        elmnt.style.top = (elmnt.offsetTop - pos2) + "px";
        elmnt.style.left = (elmnt.offsetLeft - pos1) + "px";
    }

    function closeDragElement() {
        // stop moving when mouse button is released:
        document.onmouseup = null;
        document.onmousemove = null;
    }
}
