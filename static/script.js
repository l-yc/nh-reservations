document.getElementById('createEventForm').addEventListener('submit', function(e) {
    e.preventDefault();
    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
	const location = document.getElementById('location').value;
	const converter = (value) => {
		const fakeUtcTime = new Date(`${value}Z`);
		const d = new Date(fakeUtcTime.getTime() + fakeUtcTime.getTimezoneOffset() * 60000);
		return d.toISOString();
	}

	const startTime = converter(document.getElementById('startTime').value);
	const endTime = converter(document.getElementById('endTime').value);

	fetch('/events', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ title, description, location, start_time: startTime, end_time: endTime })
    }).then(response => {
        if (response.status === 409) {
            showConflictMessage();
        } else if (response.ok) {
            response.json().then(data => {
                alert('Event created with ID: ' + data.id);
                loadEvents();
            });
        } else {
            response.text().then(text => alert('Error: ' + text));
        }
    }).catch(error => console.error('Error:', error));
});

function loadEvents(locationFilter = '') {
    fetch('/events')
        .then(response => response.json())
        .then(data => {
            const calendarEl = document.getElementById('calendar');
            const calendar = new FullCalendar.Calendar(calendarEl, {
                initialView: 'dayGridMonth',
                events: data
                    .filter(event => !locationFilter || event.location === locationFilter)
                    .map(event => ({
                        title: event.title + ' (' + event.end_time + ')',
                        start: event.start_time,
                        end: event.end_time,
                        description: event.description,
                        location: event.location
                    })),
                eventDidMount: function(info) {
                    var tooltip = new Tooltip(info.el, {
                        title: info.event.extendedProps.description + '<br>Location: ' + info.event.extendedProps.location,
                        placement: 'top',
                        trigger: 'hover',
                        container: 'body'
                    });
                }
            });
            calendar.render();
        });
}

function showConflictMessage() {
    const banner = document.getElementById('conflictMessage');
    banner.classList.remove('hidden');
}

function dismissBanner() {
    const banner = document.getElementById('conflictMessage');
    banner.classList.add('hidden');
}

// Handle location tabs for creating events
const createTabs = document.querySelectorAll('.location-tab');
createTabs.forEach(tab => {
    tab.addEventListener('click', function() {
        createTabs.forEach(t => t.classList.remove('bg-gray-500', 'text-white'));
        createTabs.forEach(t => t.classList.add('bg-gray-300'));
        this.classList.remove('bg-gray-300');
        this.classList.add('bg-gray-500', 'text-white');
        document.getElementById('location').value = this.dataset.location;
    });
});

// Handle location tabs for filtering events
const eventTabs = document.querySelectorAll('.event-location-tab');
eventTabs.forEach(tab => {
    tab.addEventListener('click', function() {
        eventTabs.forEach(t => t.classList.remove('bg-gray-500', 'text-white'));
        eventTabs.forEach(t => t.classList.add('bg-gray-300'));
        this.classList.remove('bg-gray-300');
        this.classList.add('bg-gray-500', 'text-white');
        loadEvents(this.dataset.location);
    });
});

window.onload = function() {
    loadEvents('Conference Room 1'); // Load events for the default active location
};
