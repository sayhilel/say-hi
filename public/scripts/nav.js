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
                let currIndex = index;
                let opt = document.querySelectorAll(".selector li");
                let numOpt = opt.length
                switch (event.key) {
                    case "ArrowUp":
                        opt[index].classList.remove("highlight");
                        index = (index - 1 + numOpt) % numOpt;
                        opt[index].classList.add("highlight");
                        event.preventDefault()
                        break;
                    case "ArrowDown":
                        opt[index].classList.remove("highlight");
                        index = (index + 1) % numOpt;
                        opt[index].classList.add("highlight");
                        event.preventDefault()
                        break;

                }
                if (currIndex !== index) {
                    const url = `/projects/${index}`;
                    htmx.ajax('GET', url, { target: '#project-box' });
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
        pos3 = e.clientX;
        pos4 = e.clientY;
        document.onmouseup = closeDragElement;
        document.onmousemove = elementDrag;
    }

    function elementDrag(e) {
        e = e || window;
        e.preventDefault();
        pos1 = pos3 - e.clientX;
        pos2 = pos4 - e.clientY;
        pos3 = e.clientX;
        pos4 = e.clientY;
        elmnt.style.top = (elmnt.offsetTop - pos2) + "px";
        elmnt.style.left = (elmnt.offsetLeft - pos1) + "px";
    }

    function closeDragElement() {
        document.onmouseup = null;
        document.onmousemove = null;
    }
}
