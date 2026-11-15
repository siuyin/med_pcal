I want to use pcal to generate a caledar to remind a patient to take their medicine. 
I wish to do this by placing two or three checkboxes on each date.
These represent the medicine taking times. 
How do I go about customizing pcal -- the unix calendar utility


# Adjust the text layout: use a slash to start a new line inside the calendar box
# This places three distinct check boxes on separate lines under the date number

all days in all [ ] Morning / [ ] Afternoon / [ ] Evening

or in horizontal format:

all days in all [ ] Morning   [ ] Afternoon   [ ] Evening



pcal -f medicine.txt 11 2026 4 | ps2pdf - medicine_calendar.pdf

pdfjam --booklet true --outfile medicine_calendar-book.pdf  medicine_calendar.pdf


