function returnHome(event) {
    if (event.ctrlKey && event.key === "c") {
        htmx.ajax('GET', "about-me", '#terminal-window');
    }
}


