const converter = (value) => {
    const fakeUtcTime = new Date(`${value}Z`);
    const d = new Date(fakeUtcTime.getTime() + fakeUtcTime.getTimezoneOffset() * 60000);
    return d.toISOString();
}

const invConverter = (d) => {
    d.setMinutes(d.getMinutes() - d.getTimezoneOffset());
    return d.toISOString().slice(0, 16);
}

document.getElementById('createEventForm').addEventListener('submit', function(e) {
    e.preventDefault();
    const eventId = Number.parseInt(document.getElementById('eventId').value);

    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
	const location = document.getElementById('location').value;
	const startTime = converter(document.getElementById('startTime').value);
	const endTime = converter(document.getElementById('endTime').value);

	fetch('/events', {
        method: eventId === "" ? "POST" : "DELETE",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ id: eventId, title, description, location, start_time: startTime, end_time: endTime })
    }).then(response => {
        if (response.status === 409) {
            Swal.fire({
                title: "Failed to create event",
                text: "Event conflicts with an existing event.",
                icon: "error",
            });
        } else if (response.ok) {
            response.json().then(data => {
                Swal.fire({
                    title: "Sucessfully updated event",
                    text: "",
                    icon: "success",
                });
                // default to active location
                loadEvents(document.querySelector('#event-location').value);
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
                initialView: 'timeGridWeek',
                headerToolbar: {
                    left: 'prev,next',
                    center: 'title',
                    right: 'dayGridMonth,timeGridWeek' // user can switch between the two
                },
                events: data
                    .filter(event => !locationFilter || event.location === locationFilter)
                    .map(event => ({
                        id: event.id,
                        title: event.title,
                        start: event.start_time,
                        end: event.end_time,
                        description: event.description,
                        location: event.location
                    })),
                eventDidMount: function(info) {
                    tippy(info.el, {
                        content: '[' + info.event.creator + '] ' + info.event.extendedProps.description,
                    });
                },
                eventClick: function(info) {
                    // fill in the createEventForm with event details
                    document.getElementById('eventId').value = info.event.id;
                    document.getElementById('title').value = info.event.title;
                    document.getElementById('description').value = info.event.extendedProps.description;
                    document.getElementById('location').value = info.event.extendedProps.location;

                    document.getElementById('startTime').value = invConverter(info.event.start);
                    document.getElementById('endTime').value = invConverter(info.event.end);
                    setReservationForm('delete');
                }
            });
            calendar.render();
        });
}

function setReservationForm(action) {
    const setCreateButton = document.querySelector('#setCreateButton');
    const submitButton = document.getElementById('submitButton');
    const title = document.getElementById('formTitle');
    if (action === 'create') {
        setCreateButton.classList.add('hidden');
        submitButton.innerText = 'Create';
        title.innerText = 'Create Reservation';
        document.getElementById('eventId').value = '';
    } else if (action === 'delete') {
        setCreateButton.classList.remove('hidden');
        submitButton.innerText = 'Delete';
        title.innerText = 'Delete Reservation';
    }
}

const setCreateButton = document.querySelector('#setCreateButton');
setCreateButton.addEventListener('click', function() {
    setReservationForm('create');
});

// Handle location tabs for creating events
const createTabs = document.querySelectorAll('.location-tab');
createTabs.forEach(tab => {
    tab.addEventListener('click', function() {
        createTabs.forEach(t => t.classList.remove('bg-gray-500', 'text-white'));
        createTabs.forEach(t => t.classList.add('bg-gray-300'));
        tab.classList.remove('bg-gray-300');
        tab.classList.add('bg-gray-500', 'text-white');
        document.getElementById('location').value = tab.dataset.location;
    });
});

// Handle location tabs for filtering events
const eventTabs = document.querySelectorAll('.event-location-tab');
eventTabs.forEach(tab => {
    tab.addEventListener('click', function() {
        eventTabs.forEach(t => t.classList.remove('bg-gray-500', 'text-white'));
        eventTabs.forEach(t => t.classList.add('bg-gray-300'));
        tab.classList.remove('bg-gray-300');
        tab.classList.add('bg-gray-500', 'text-white');
        loadEvents(tab.dataset.location);
        document.getElementById('event-location').value = tab.dataset.location;
    });
});

window.onload = function() {
    loadEvents('MPR'); // Load events for the default active location
};
