import 'package:flutter/material.dart';

class AppLocalizations {
  final Locale locale;

  AppLocalizations(this.locale);

  static AppLocalizations? of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations);
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  static final List<Locale> supportedLocales = [
    const Locale('en', ''), // English
    const Locale('id', ''), // Indonesian
  ];

  // Profile Screen
  String get profile =>
      _localizedValues[locale.languageCode]?['profile'] ?? 'Profile';
  String get settings =>
      _localizedValues[locale.languageCode]?['settings'] ?? 'Settings';
  String get notifications =>
      _localizedValues[locale.languageCode]?['notifications'] ??
      'Notifications';
  String get manageNotificationSettings =>
      _localizedValues[locale.languageCode]?['manageNotificationSettings'] ??
      'Manage notification settings';
  String get noNotificationsFound =>
      _localizedValues[locale.languageCode]?['noNotificationsFound'] ??
      'No notifications found';
  String get markAllAsRead =>
      _localizedValues[locale.languageCode]?['markAllAsRead'] ??
      'Mark All as Read';
  String get markAsRead =>
      _localizedValues[locale.languageCode]?['markAsRead'] ?? 'Mark as Read';
  String get filter =>
      _localizedValues[locale.languageCode]?['filter'] ?? 'Filter';
  String get unread =>
      _localizedValues[locale.languageCode]?['unread'] ?? 'Unread';
  String get read => _localizedValues[locale.languageCode]?['read'] ?? 'Read';
  String get filterNotifications =>
      _localizedValues[locale.languageCode]?['filterNotifications'] ??
      'Filter Notifications';
  String get notificationMarkedAsRead =>
      _localizedValues[locale.languageCode]?['notificationMarkedAsRead'] ??
      'Notification marked as read';
  String get failedToMarkAsRead =>
      _localizedValues[locale.languageCode]?['failedToMarkAsRead'] ??
      'Failed to mark as read';
  String get allNotificationsMarkedAsRead =>
      _localizedValues[locale.languageCode]?['allNotificationsMarkedAsRead'] ??
      'All notifications marked as read';
  String get failedToMarkAllAsRead =>
      _localizedValues[locale.languageCode]?['failedToMarkAllAsRead'] ??
      'Failed to mark all as read';
  String get deleteNotification =>
      _localizedValues[locale.languageCode]?['deleteNotification'] ??
      'Delete Notification';
  String get deleteNotificationConfirmation =>
      _localizedValues[locale
          .languageCode]?['deleteNotificationConfirmation'] ??
      'Are you sure you want to delete this notification?';
  String get notificationDeleted =>
      _localizedValues[locale.languageCode]?['notificationDeleted'] ??
      'Notification deleted successfully';
  String get failedToDeleteNotification =>
      _localizedValues[locale.languageCode]?['failedToDeleteNotification'] ??
      'Failed to delete notification';
  String get language =>
      _localizedValues[locale.languageCode]?['language'] ?? 'Language';
  String get theme =>
      _localizedValues[locale.languageCode]?['theme'] ?? 'Theme';
  String get lightTheme =>
      _localizedValues[locale.languageCode]?['lightTheme'] ?? 'Light theme';
  String get darkTheme =>
      _localizedValues[locale.languageCode]?['darkTheme'] ?? 'Dark theme';
  String get systemTheme =>
      _localizedValues[locale.languageCode]?['systemTheme'] ??
      'Follow system setting';
  String get about =>
      _localizedValues[locale.languageCode]?['about'] ?? 'About';
  String get appVersion =>
      _localizedValues[locale.languageCode]?['appVersion'] ?? 'App Version';
  String get privacyPolicy =>
      _localizedValues[locale.languageCode]?['privacyPolicy'] ??
      'Privacy Policy';
  String get viewPrivacyPolicy =>
      _localizedValues[locale.languageCode]?['viewPrivacyPolicy'] ??
      'View privacy policy';
  String get termsOfService =>
      _localizedValues[locale.languageCode]?['termsOfService'] ??
      'Terms of Service';
  String get viewTermsOfService =>
      _localizedValues[locale.languageCode]?['viewTermsOfService'] ??
      'View terms of service';
  String get firstName =>
      _localizedValues[locale.languageCode]?['firstName'] ?? 'First Name';
  String get lastName =>
      _localizedValues[locale.languageCode]?['lastName'] ?? 'Last Name';
  String get logout =>
      _localizedValues[locale.languageCode]?['logout'] ?? 'Logout';
  String get logoutConfirmation =>
      _localizedValues[locale.languageCode]?['logoutConfirmation'] ??
      'Are you sure you want to logout?';
  String get cancel =>
      _localizedValues[locale.languageCode]?['cancel'] ?? 'Cancel';
  String get selectTheme =>
      _localizedValues[locale.languageCode]?['selectTheme'] ?? 'Select Theme';
  String get selectLanguage =>
      _localizedValues[locale.languageCode]?['selectLanguage'] ??
      'Select Language';
  String get english =>
      _localizedValues[locale.languageCode]?['english'] ?? 'English';
  String get indonesian =>
      _localizedValues[locale.languageCode]?['indonesian'] ?? 'Indonesian';

  // Login Screen
  String get signInToYourAccount =>
      _localizedValues[locale.languageCode]?['signInToYourAccount'] ??
      'Sign in to your Account';
  String get enterEmailPassword =>
      _localizedValues[locale.languageCode]?['enterEmailPassword'] ??
      'Enter your email and password to log in';
  String get email =>
      _localizedValues[locale.languageCode]?['email'] ?? 'Email';
  String get password =>
      _localizedValues[locale.languageCode]?['password'] ?? 'Password';
  String get enterPassword =>
      _localizedValues[locale.languageCode]?['enterPassword'] ??
      'Enter your password';
  String get rememberMe =>
      _localizedValues[locale.languageCode]?['rememberMe'] ?? 'Remember me';
  String get logIn =>
      _localizedValues[locale.languageCode]?['logIn'] ?? 'Log In';

  // Bottom Navigation
  String get home => _localizedValues[locale.languageCode]?['home'] ?? 'Home';
  String get accounts =>
      _localizedValues[locale.languageCode]?['accounts'] ?? 'Accounts';
  String get contacts =>
      _localizedValues[locale.languageCode]?['contacts'] ?? 'Contacts';
  String get accountsAndContacts =>
      _localizedValues[locale.languageCode]?['accountsAndContacts'] ??
      'Accounts';
  String get reports =>
      _localizedValues[locale.languageCode]?['reports'] ?? 'Reports';
  String get tasks =>
      _localizedValues[locale.languageCode]?['tasks'] ?? 'Tasks';
  String get reportsAndTasks =>
      _localizedValues[locale.languageCode]?['reportsAndTasks'] ?? 'Reports';
  String get visitsMenu =>
      _localizedValues[locale.languageCode]?['visitsMenu'] ?? 'Visits';
  String get searchVisitReports =>
      _localizedValues[locale.languageCode]?['searchVisitReports'] ??
      'Search visit reports...';
  String get search =>
      _localizedValues[locale.languageCode]?['search'] ?? 'Search';
  String get searchTasks =>
      _localizedValues[locale.languageCode]?['searchTasks'] ??
      'Search tasks...';
  String get searchSchedules =>
      _localizedValues[locale.languageCode]?['searchSchedules'] ??
      'Search schedules...';
  String get noSchedulesFound =>
      _localizedValues[locale.languageCode]?['noSchedulesFound'] ??
      'No schedules found';
  String get tapToCreateSchedule =>
      _localizedValues[locale.languageCode]?['tapToCreateSchedule'] ??
      'Tap + to create a new schedule';
  String get schedules =>
      _localizedValues[locale.languageCode]?['schedules'] ?? 'Schedules';
  String get scheduleDetails =>
      _localizedValues[locale.languageCode]?['scheduleDetails'] ??
      'Schedule Details';
  String get scheduleInformation =>
      _localizedValues[locale.languageCode]?['scheduleInformation'] ??
      'Schedule Information';
  String get scheduledAt =>
      _localizedValues[locale.languageCode]?['scheduledAt'] ?? 'Scheduled At';
  String get reminderMinutesBefore =>
      _localizedValues[locale.languageCode]?['reminderMinutesBefore'] ??
      'Reminder (minutes before)';
  String get invalidReminderMinutes =>
      _localizedValues[locale.languageCode]?['invalidReminderMinutes'] ??
      'Invalid reminder minutes (0-10080)';
  String get enterTitle =>
      _localizedValues[locale.languageCode]?['enterTitle'] ?? 'Enter title';
  String get reminder =>
      _localizedValues[locale.languageCode]?['reminder'] ?? 'Reminder';
  String get minutesBefore =>
      _localizedValues[locale.languageCode]?['minutesBefore'] ??
      'minutes before';
  String get linkedTask =>
      _localizedValues[locale.languageCode]?['linkedTask'] ?? 'Linked Task';
  String get selectTask =>
      _localizedValues[locale.languageCode]?['selectTask'] ?? 'Select Task';
  String get scheduledAtRequired =>
      _localizedValues[locale.languageCode]?['scheduledAtRequired'] ??
      'Scheduled time is required';
  String get scheduledFrom =>
      _localizedValues[locale.languageCode]?['scheduledFrom'] ??
      'Scheduled From';
  String get scheduledTo =>
      _localizedValues[locale.languageCode]?['scheduledTo'] ?? 'Scheduled To';
  String get selectDate =>
      _localizedValues[locale.languageCode]?['selectDate'] ?? 'Select Date';
  String get selectScheduledFrom =>
      _localizedValues[locale.languageCode]?['selectScheduledFrom'] ??
      'Select scheduled from';
  String get selectScheduledTo =>
      _localizedValues[locale.languageCode]?['selectScheduledTo'] ??
      'Select scheduled to';
  String get clearDate =>
      _localizedValues[locale.languageCode]?['clearDate'] ?? 'Clear Date';
  String get noVisitReportsFound =>
      _localizedValues[locale.languageCode]?['noVisitReportsFound'] ??
      'No visit reports found';
  String get noData =>
      _localizedValues[locale.languageCode]?['noData'] ?? 'No data';
  String get noDescription =>
      _localizedValues[locale.languageCode]?['noDescription'] ??
      'No description';
  String get location =>
      _localizedValues[locale.languageCode]?['location'] ?? 'Location';
  String get noLocation =>
      _localizedValues[locale.languageCode]?['noLocation'] ?? 'No location';
  String get relatedVisitReport =>
      _localizedValues[locale.languageCode]?['relatedVisitReport'] ??
      'Related Visit Report';
  String get viewVisitReport =>
      _localizedValues[locale.languageCode]?['viewVisitReport'] ??
      'View Visit Report';
  String get activityType =>
      _localizedValues[locale.languageCode]?['activityType'] ?? 'Activity Type';
  String get scheduled =>
      _localizedValues[locale.languageCode]?['scheduled'] ?? 'Scheduled';
  String get confirmed =>
      _localizedValues[locale.languageCode]?['confirmed'] ?? 'Confirmed';
  String get createSchedule =>
      _localizedValues[locale.languageCode]?['createSchedule'] ??
      'Create Schedule';
  String get editSchedule =>
      _localizedValues[locale.languageCode]?['editSchedule'] ?? 'Edit Schedule';
  String get updateSchedule =>
      _localizedValues[locale.languageCode]?['updateSchedule'] ??
      'Update Schedule';
  String get deleteSchedule =>
      _localizedValues[locale.languageCode]?['deleteSchedule'] ??
      'Delete Schedule';
  String get deleteScheduleConfirmation =>
      _localizedValues[locale.languageCode]?['deleteScheduleConfirmation'] ??
      'Are you sure you want to delete this schedule? This action cannot be undone.';
  String get scheduleCreatedSuccessfully =>
      _localizedValues[locale.languageCode]?['scheduleCreatedSuccessfully'] ??
      'Schedule created successfully';
  String get scheduleUpdatedSuccessfully =>
      _localizedValues[locale.languageCode]?['scheduleUpdatedSuccessfully'] ??
      'Schedule updated successfully';
  String get scheduleDeletedSuccessfully =>
      _localizedValues[locale.languageCode]?['scheduleDeletedSuccessfully'] ??
      'Schedule deleted successfully';
  String get failedToCreateSchedule =>
      _localizedValues[locale.languageCode]?['failedToCreateSchedule'] ??
      'Failed to create schedule';
  String get failedToUpdateSchedule =>
      _localizedValues[locale.languageCode]?['failedToUpdateSchedule'] ??
      'Failed to update schedule';
  String get failedToDeleteSchedule =>
      _localizedValues[locale.languageCode]?['failedToDeleteSchedule'] ??
      'Failed to delete schedule';
  String get noTasksFound =>
      _localizedValues[locale.languageCode]?['noTasksFound'] ??
      'No tasks found';
  String get noSearchResults =>
      _localizedValues[locale.languageCode]?['noSearchResults'] ??
      'No results found for';
  String get searchResultsFor =>
      _localizedValues[locale.languageCode]?['searchResultsFor'] ??
      'Search results for';
  String get tapToCreateVisitReport =>
      _localizedValues[locale.languageCode]?['tapToCreateVisitReport'] ??
      'Tap + to create a new visit report';
  String get tapToCreateTask =>
      _localizedValues[locale.languageCode]?['tapToCreateTask'] ??
      'Tap + to create a new task';
  String get createVisitReport =>
      _localizedValues[locale.languageCode]?['createVisitReport'] ??
      'Create Visit Report';
  String get visitReports =>
      _localizedValues[locale.languageCode]?['visitReports'] ?? 'Visit Reports';
  String get visitReportDetails =>
      _localizedValues[locale.languageCode]?['visitReportDetails'] ??
      'Visit Report Details';
  String get visitInformation =>
      _localizedValues[locale.languageCode]?['visitInformation'] ??
      'Visit Information';
  String get checkInOutStatus =>
      _localizedValues[locale.languageCode]?['checkInOutStatus'] ??
      'Check-in/out Status';
  String get visitDate =>
      _localizedValues[locale.languageCode]?['visitDate'] ?? 'Visit Date';
  String get purpose =>
      _localizedValues[locale.languageCode]?['purpose'] ?? 'Purpose';
  String get photos =>
      _localizedValues[locale.languageCode]?['photos'] ?? 'Photos';
  String get checkInTime =>
      _localizedValues[locale.languageCode]?['checkInTime'] ?? 'Check-in Time';
  String get checkOutTime =>
      _localizedValues[locale.languageCode]?['checkOutTime'] ??
      'Check-out Time';
  String get checkInLocation =>
      _localizedValues[locale.languageCode]?['checkInLocation'] ??
      'Check-in Location';
  String get checkOutLocation =>
      _localizedValues[locale.languageCode]?['checkOutLocation'] ??
      'Check-out Location';
  String get checkIn =>
      _localizedValues[locale.languageCode]?['checkIn'] ?? 'Check In';
  String get checkOut =>
      _localizedValues[locale.languageCode]?['checkOut'] ?? 'Check Out';
  String get uploadPhoto =>
      _localizedValues[locale.languageCode]?['uploadPhoto'] ?? 'Upload Photo';
  String get notCheckedIn =>
      _localizedValues[locale.languageCode]?['notCheckedIn'] ??
      'Not checked in';
  String get notCheckedOut =>
      _localizedValues[locale.languageCode]?['notCheckedOut'] ??
      'Not checked out';
  String get checkInSuccessful =>
      _localizedValues[locale.languageCode]?['checkInSuccessful'] ??
      'Check-in successful';
  String get checkOutSuccessful =>
      _localizedValues[locale.languageCode]?['checkOutSuccessful'] ??
      'Check-out successful';
  String get photoUploadedSuccessfully =>
      _localizedValues[locale.languageCode]?['photoUploadedSuccessfully'] ??
      'Photo uploaded successfully';
  String get failedToCheckIn =>
      _localizedValues[locale.languageCode]?['failedToCheckIn'] ??
      'Failed to check in';
  String get failedToCheckOut =>
      _localizedValues[locale.languageCode]?['failedToCheckOut'] ??
      'Failed to check out';
  String get failedToUploadPhoto =>
      _localizedValues[locale.languageCode]?['failedToUploadPhoto'] ??
      'Failed to upload photo';
  String get selfieRequiredForCheckIn =>
      _localizedValues[locale.languageCode]?['selfieRequiredForCheckIn'] ??
      'Selfie picture is required for check-in. Please take a photo to continue.';
  String get previewSelfie =>
      _localizedValues[locale.languageCode]?['previewSelfie'] ??
      'Preview Selfie';
  String get retake =>
      _localizedValues[locale.languageCode]?['retake'] ?? 'Retake';
  String get confirmCheckOut =>
      _localizedValues[locale.languageCode]?['confirmCheckOut'] ??
      'Confirm Check-out';
  String get checkOutLocationRequired =>
      _localizedValues[locale.languageCode]?['checkOutLocationRequired'] ??
      'Current location is required for check-out';
  String get gettingLocation =>
      _localizedValues[locale.languageCode]?['gettingLocation'] ??
      'Getting your current location...';
  String get locationAccuracy =>
      _localizedValues[locale.languageCode]?['locationAccuracy'] ??
      'Location Accuracy';
  String get currentLocation =>
      _localizedValues[locale.languageCode]?['currentLocation'] ??
      'Current Location';
  String get fakeGPSDetected =>
      _localizedValues[locale.languageCode]?['fakeGPSDetected'] ??
      'Fake GPS Detected';
  String get fakeGPSDescription =>
      _localizedValues[locale.languageCode]?['fakeGPSDescription'] ??
      'We detected that you are using a Fake GPS application. Check-in and check-out features are disabled when Fake GPS is active.';
  String get detectedReason =>
      _localizedValues[locale.languageCode]?['detectedReason'] ??
      'Detected Reason';
  String get howToDisableFakeGPS =>
      _localizedValues[locale.languageCode]?['howToDisableFakeGPS'] ??
      'How to Disable Fake GPS';
  String get fakeGPSStep1 =>
      _localizedValues[locale.languageCode]?['fakeGPSStep1'] ??
      'Close or uninstall any Fake GPS or location spoofing applications on your device';
  String get fakeGPSStep2 =>
      _localizedValues[locale.languageCode]?['fakeGPSStep2'] ??
      'Restart your device to ensure all location services are reset';
  String get fakeGPSStep3 =>
      _localizedValues[locale.languageCode]?['fakeGPSStep3'] ??
      'Clear app cache and location permissions';
  String get fakeGPSStep4 =>
      _localizedValues[locale.languageCode]?['fakeGPSStep4'] ??
      'Try check-in/check-out again after disabling Fake GPS';
  String get fakeGPSImportantNote =>
      _localizedValues[locale.languageCode]?['fakeGPSImportantNote'] ??
      'Check-in and check-out features require real GPS location for verification. Using Fake GPS is not allowed and will prevent you from using these features.';
  String get close =>
      _localizedValues[locale.languageCode]?['close'] ?? 'Close';
  String get selectContact =>
      _localizedValues[locale.languageCode]?['selectContact'] ??
      'Select Contact';
  String get selectContactOptional =>
      _localizedValues[locale.languageCode]?['selectContactOptional'] ??
      'Select contact (optional)';
  String get selectVisitDate =>
      _localizedValues[locale.languageCode]?['selectVisitDate'] ??
      'Select visit date';
  String get enterVisitPurpose =>
      _localizedValues[locale.languageCode]?['enterVisitPurpose'] ??
      'Enter visit purpose';
  String get enterAdditionalNotes =>
      _localizedValues[locale.languageCode]?['enterAdditionalNotes'] ??
      'Enter additional notes';
  String get pleaseSelectAccount =>
      _localizedValues[locale.languageCode]?['pleaseSelectAccount'] ??
      'Please select an account';
  String get visitReportCreatedSuccessfully =>
      _localizedValues[locale
          .languageCode]?['visitReportCreatedSuccessfully'] ??
      'Visit report created successfully';
  String get failedToCreateVisitReport =>
      _localizedValues[locale.languageCode]?['failedToCreateVisitReport'] ??
      'Failed to create visit report';
  String get deal => _localizedValues[locale.languageCode]?['deal'] ?? 'Deal';
  String get lead => _localizedValues[locale.languageCode]?['lead'] ?? 'Lead';
  String get none => _localizedValues[locale.languageCode]?['none'] ?? 'None';
  String get pleaseSelectDeal =>
      _localizedValues[locale.languageCode]?['pleaseSelectDeal'] ??
      'Please select a deal';
  String get pleaseSelectLead =>
      _localizedValues[locale.languageCode]?['pleaseSelectLead'] ??
      'Please select a lead';
  String get noLeadsAvailable =>
      _localizedValues[locale.languageCode]?['noLeadsAvailable'] ??
      'No leads available';
  String get purposeRequired =>
      _localizedValues[locale.languageCode]?['purposeRequired'] ??
      'Purpose is required';
  String get updateVisitReport =>
      _localizedValues[locale.languageCode]?['updateVisitReport'] ??
      'Update Visit Report';
  String get deleteVisitReport =>
      _localizedValues[locale.languageCode]?['deleteVisitReport'] ??
      'Delete Visit Report';
  String get submitVisitReport =>
      _localizedValues[locale.languageCode]?['submitVisitReport'] ??
      'Submit Visit Report';
  String get deleteVisitReportConfirmation =>
      _localizedValues[locale.languageCode]?['deleteVisitReportConfirmation'] ??
      'Are you sure you want to delete this visit report? This action cannot be undone.';
  String get visitReportUpdatedSuccessfully =>
      _localizedValues[locale
          .languageCode]?['visitReportUpdatedSuccessfully'] ??
      'Visit report updated successfully';
  String get visitReportDeletedSuccessfully =>
      _localizedValues[locale
          .languageCode]?['visitReportDeletedSuccessfully'] ??
      'Visit report deleted successfully';
  String get visitReportSubmittedSuccessfully =>
      _localizedValues[locale
          .languageCode]?['visitReportSubmittedSuccessfully'] ??
      'Visit report submitted successfully';
  String get failedToUpdateVisitReport =>
      _localizedValues[locale.languageCode]?['failedToUpdateVisitReport'] ??
      'Failed to update visit report';
  String get failedToDeleteVisitReport =>
      _localizedValues[locale.languageCode]?['failedToDeleteVisitReport'] ??
      'Failed to delete visit report';
  String get failedToSubmitVisitReport =>
      _localizedValues[locale.languageCode]?['failedToSubmitVisitReport'] ??
      'Failed to submit visit report';
  String get outcome =>
      _localizedValues[locale.languageCode]?['outcome'] ?? 'Outcome';
  String get nextSteps =>
      _localizedValues[locale.languageCode]?['nextSteps'] ?? 'Next Steps';
  String get selectOutcome =>
      _localizedValues[locale.languageCode]?['selectOutcome'] ??
      'Select outcome';
  String get enterNextSteps =>
      _localizedValues[locale.languageCode]?['enterNextSteps'] ??
      'Enter next steps (optional)';
  String get all => _localizedValues[locale.languageCode]?['all'] ?? 'All';
  String get filterByStatus =>
      _localizedValues[locale.languageCode]?['filterByStatus'] ??
      'Filter by Status';
  String get filterByScheduledDate =>
      _localizedValues[locale.languageCode]?['filterByScheduledDate'] ??
      'Filter by Scheduled Date';
  String get filterByPriority =>
      _localizedValues[locale.languageCode]?['filterByPriority'] ??
      'Filter by Priority';
  String get filterByType =>
      _localizedValues[locale.languageCode]?['filterByType'] ??
      'Filter by Type';
  String get filterByDueDate =>
      _localizedValues[locale.languageCode]?['filterByDueDate'] ??
      'Filter by Due Date';
  String get dueDateFrom =>
      _localizedValues[locale.languageCode]?['dueDateFrom'] ?? 'Due Date From';
  String get dueDateTo =>
      _localizedValues[locale.languageCode]?['dueDateTo'] ?? 'Due Date To';
  String get selectDueDateFrom =>
      _localizedValues[locale.languageCode]?['selectDueDateFrom'] ??
      'Select due date from';
  String get selectDueDateTo =>
      _localizedValues[locale.languageCode]?['selectDueDateTo'] ??
      'Select due date to';
  String get clearFilters =>
      _localizedValues[locale.languageCode]?['clearFilters'] ?? 'Clear filters';
  String get clearSearch =>
      _localizedValues[locale.languageCode]?['clearSearch'] ?? 'Clear Search';
  String get taskDetails =>
      _localizedValues[locale.languageCode]?['taskDetails'] ?? 'Task Details';
  String get taskInformation =>
      _localizedValues[locale.languageCode]?['taskInformation'] ??
      'Task Information';
  String get relatedInformation =>
      _localizedValues[locale.languageCode]?['relatedInformation'] ??
      'Related Information';
  String get reminders =>
      _localizedValues[locale.languageCode]?['reminders'] ?? 'Reminders';
  String get completeTask =>
      _localizedValues[locale.languageCode]?['completeTask'] ?? 'Complete Task';
  String get markInProgress =>
      _localizedValues[locale.languageCode]?['markInProgress'] ??
      'Mark In Progress';
  String get markInProgressConfirmation =>
      _localizedValues[locale.languageCode]?['markInProgressConfirmation'] ??
      'Are you sure you want to mark this task as in progress?';
  String get taskMarkedInProgress =>
      _localizedValues[locale.languageCode]?['taskMarkedInProgress'] ??
      'Task marked as in progress';
  String get failedToMarkInProgress =>
      _localizedValues[locale.languageCode]?['failedToMarkInProgress'] ??
      'Failed to mark task as in progress';
  String get addReminder =>
      _localizedValues[locale.languageCode]?['addReminder'] ?? 'Add Reminder';
  String get completeTaskConfirmation =>
      _localizedValues[locale.languageCode]?['completeTaskConfirmation'] ??
      'Are you sure you want to mark this task as completed?';
  String get taskCompletedSuccessfully =>
      _localizedValues[locale.languageCode]?['taskCompletedSuccessfully'] ??
      'Task completed successfully';
  String get failedToCompleteTask =>
      _localizedValues[locale.languageCode]?['failedToCompleteTask'] ??
      'Failed to complete task';
  String get deleteTask =>
      _localizedValues[locale.languageCode]?['deleteTask'] ?? 'Delete Task';
  String get deleteTaskConfirmation =>
      _localizedValues[locale.languageCode]?['deleteTaskConfirmation'] ??
      'Are you sure you want to delete this task? This action cannot be undone.';
  String get taskDeletedSuccessfully =>
      _localizedValues[locale.languageCode]?['taskDeletedSuccessfully'] ??
      'Task deleted successfully';
  String get failedToDeleteTask =>
      _localizedValues[locale.languageCode]?['failedToDeleteTask'] ??
      'Failed to delete task';
  String get reminderCreatedSuccessfully =>
      _localizedValues[locale.languageCode]?['reminderCreatedSuccessfully'] ??
      'Reminder created successfully';
  String get failedToCreateReminder =>
      _localizedValues[locale.languageCode]?['failedToCreateReminder'] ??
      'Failed to create reminder';
  String get selectReminderDate =>
      _localizedValues[locale.languageCode]?['selectReminderDate'] ??
      'Select reminder date and time';
  String get reminderMessage =>
      _localizedValues[locale.languageCode]?['reminderMessage'] ??
      'Message (optional)';
  String get enterReminderMessage =>
      _localizedValues[locale.languageCode]?['enterReminderMessage'] ??
      'Enter reminder message';
  String get title =>
      _localizedValues[locale.languageCode]?['title'] ?? 'Title';
  String get description =>
      _localizedValues[locale.languageCode]?['description'] ?? 'Description';
  String get type => _localizedValues[locale.languageCode]?['type'] ?? 'Type';
  String get dueDate =>
      _localizedValues[locale.languageCode]?['dueDate'] ?? 'Due Date';
  String get selectDueDate =>
      _localizedValues[locale.languageCode]?['selectDueDate'] ??
      'Select due date';
  String get createTask =>
      _localizedValues[locale.languageCode]?['createTask'] ?? 'Create Task';
  String get editTask =>
      _localizedValues[locale.languageCode]?['editTask'] ?? 'Edit Task';
  String get updateTask =>
      _localizedValues[locale.languageCode]?['updateTask'] ?? 'Update Task';
  String get enterTaskTitle =>
      _localizedValues[locale.languageCode]?['enterTaskTitle'] ??
      'Enter task title';
  String get enterTaskDescription =>
      _localizedValues[locale.languageCode]?['enterTaskDescription'] ??
      'Enter task description';
  String get taskCreatedSuccessfully =>
      _localizedValues[locale.languageCode]?['taskCreatedSuccessfully'] ??
      'Task created successfully';
  String get taskUpdatedSuccessfully =>
      _localizedValues[locale.languageCode]?['taskUpdatedSuccessfully'] ??
      'Task updated successfully';
  String get failedToCreateTask =>
      _localizedValues[locale.languageCode]?['failedToCreateTask'] ??
      'Failed to create task';
  String get failedToUpdateTask =>
      _localizedValues[locale.languageCode]?['failedToUpdateTask'] ??
      'Failed to update task';
  String get titleIsRequired =>
      _localizedValues[locale.languageCode]?['titleIsRequired'] ??
      'Title is required';
  String get sent => _localizedValues[locale.languageCode]?['sent'] ?? 'Sent';

  // Dashboard
  String get dashboard =>
      _localizedValues[locale.languageCode]?['dashboard'] ?? 'Dashboard';
  String get welcomeBack =>
      _localizedValues[locale.languageCode]?['welcomeBack'] ?? 'Welcome back,';
  String get today =>
      _localizedValues[locale.languageCode]?['today'] ?? 'Today';
  String get thisWeek =>
      _localizedValues[locale.languageCode]?['thisWeek'] ?? 'This Week';
  String get thisMonth =>
      _localizedValues[locale.languageCode]?['thisMonth'] ?? 'This Month';
  String get thisYear =>
      _localizedValues[locale.languageCode]?['thisYear'] ?? 'This Year';
  String get totalVisits =>
      _localizedValues[locale.languageCode]?['totalVisits'] ?? 'Total Visits';
  String get totalAccounts =>
      _localizedValues[locale.languageCode]?['totalAccounts'] ??
      'Total Accounts';
  String get totalActivities =>
      _localizedValues[locale.languageCode]?['totalActivities'] ??
      'Total Activities';
  String get revenue =>
      _localizedValues[locale.languageCode]?['revenue'] ?? 'Revenue';
  String get completed =>
      _localizedValues[locale.languageCode]?['completed'] ?? 'completed';
  String get planned =>
      _localizedValues[locale.languageCode]?['planned'] ?? 'Planned';
  String get cancelled =>
      _localizedValues[locale.languageCode]?['cancelled'] ?? 'Cancelled';
  String get seeAll =>
      _localizedValues[locale.languageCode]?['seeAll'] ?? 'See All';
  String get noVisitsFound =>
      _localizedValues[locale.languageCode]?['noVisitsFound'] ??
      'No visits found';
  String get remaining =>
      _localizedValues[locale.languageCode]?['remaining'] ?? 'remaining';
  String get tersisa =>
      _localizedValues[locale.languageCode]?['tersisa'] ?? 'tersisa';
  String get active =>
      _localizedValues[locale.languageCode]?['active'] ?? 'active';
  String get visits =>
      _localizedValues[locale.languageCode]?['visits'] ?? 'visits';
  String get calls =>
      _localizedValues[locale.languageCode]?['calls'] ?? 'calls';
  String get fromWonDeals =>
      _localizedValues[locale.languageCode]?['fromWonDeals'] ??
      'From won deals';
  String get salesTarget =>
      _localizedValues[locale.languageCode]?['salesTarget'] ?? 'Sales Target';
  String get totalDeals =>
      _localizedValues[locale.languageCode]?['totalDeals'] ?? 'Total Deals';
  String get leadsBySource =>
      _localizedValues[locale.languageCode]?['leadsBySource'] ??
      'Leads by Source';
  String get upcomingTasks =>
      _localizedValues[locale.languageCode]?['upcomingTasks'] ??
      'Upcoming Tasks';
  String get pipelineSummary =>
      _localizedValues[locale.languageCode]?['pipelineSummary'] ??
      'Pipeline Summary';
  String get salesPipeline =>
      _localizedValues[locale.languageCode]?['salesPipeline'] ??
      'Sales Pipeline';
  String get recentActivities =>
      _localizedValues[locale.languageCode]?['recentActivities'] ??
      'Recent Activities';
  String get target =>
      _localizedValues[locale.languageCode]?['target'] ?? 'Target';
  String get achieved =>
      _localizedValues[locale.languageCode]?['achieved'] ?? 'Achieved';
  String get open => _localizedValues[locale.languageCode]?['open'] ?? 'Open';
  String get won => _localizedValues[locale.languageCode]?['won'] ?? 'Won';
  String get lost => _localizedValues[locale.languageCode]?['lost'] ?? 'Lost';
  String get totalValue =>
      _localizedValues[locale.languageCode]?['totalValue'] ?? 'Total Value';
  String get totalLeads =>
      _localizedValues[locale.languageCode]?['totalLeads'] ?? 'total leads';
  String get noLeadsForPeriod =>
      _localizedValues[locale.languageCode]?['noLeadsForPeriod'] ??
      'No leads for this period';
  String get noUpcomingTasks =>
      _localizedValues[locale.languageCode]?['noUpcomingTasks'] ??
      'No upcoming tasks';
  String targetProgressDescription(double progress) =>
      _localizedValues[locale.languageCode]?['targetProgressDescription']
          ?.replaceAll('{progress}', progress.toStringAsFixed(0)) ??
      '${progress.toStringAsFixed(0)}% of target achieved';
  String totalAccountsDescription(int active, int inactive) =>
      _localizedValues[locale.languageCode]?['totalAccountsDescription']
          ?.replaceAll('{active}', active.toString())
          .replaceAll('{inactive}', inactive.toString()) ??
      '$active active, $inactive inactive';
  String totalDealsDescription(int open, int won) =>
      _localizedValues[locale.languageCode]?['totalDealsDescription']
          ?.replaceAll('{open}', open.toString())
          .replaceAll('{won}', won.toString()) ??
      '$open open, $won won';
  String get totalRevenueDescription =>
      _localizedValues[locale.languageCode]?['totalRevenueDescription'] ??
      'Based on won deals in this period';
  String get topAccounts =>
      _localizedValues[locale.languageCode]?['topAccounts'] ?? 'Top Accounts';
  String get topSalesReps =>
      _localizedValues[locale.languageCode]?['topSalesReps'] ??
      'Top Sales Reps';
  String get visitStatistics =>
      _localizedValues[locale.languageCode]?['visitStatistics'] ??
      'Visit Statistics';
  String get activityTrends =>
      _localizedValues[locale.languageCode]?['activityTrends'] ??
      'Activity Trends';
  String get noTopAccounts =>
      _localizedValues[locale.languageCode]?['noTopAccounts'] ??
      'No accounts found';
  String get noTopSalesReps =>
      _localizedValues[locale.languageCode]?['noTopSalesReps'] ??
      'No sales reps found';
  String topAccountsVisits(int count) =>
      _localizedValues[locale.languageCode]?['topAccountsVisits']?.replaceAll(
        '{count}',
        count.toString(),
      ) ??
      '$count visits';
  String topAccountsActivities(int count) =>
      _localizedValues[locale.languageCode]?['topAccountsActivities']
          ?.replaceAll('{count}', count.toString()) ??
      '$count activities';
  String topSalesRepsVisits(int count) =>
      _localizedValues[locale.languageCode]?['topSalesRepsVisits']?.replaceAll(
        '{count}',
        count.toString(),
      ) ??
      '$count visits';
  String topSalesRepsAccounts(int count) =>
      _localizedValues[locale.languageCode]?['topSalesRepsAccounts']
          ?.replaceAll('{count}', count.toString()) ??
      '$count accounts';
  String get emails =>
      _localizedValues[locale.languageCode]?['emails'] ?? 'emails';
  String get pending =>
      _localizedValues[locale.languageCode]?['pending'] ?? 'pending';
  String get approved =>
      _localizedValues[locale.languageCode]?['approved'] ?? 'approved';
  String get draft =>
      _localizedValues[locale.languageCode]?['draft'] ?? 'Draft';
  String get submitted =>
      _localizedValues[locale.languageCode]?['submitted'] ?? 'Submitted';
  String get rejected =>
      _localizedValues[locale.languageCode]?['rejected'] ?? 'Rejected';
  String get total =>
      _localizedValues[locale.languageCode]?['total'] ?? 'Total';
  String get visitsToday =>
      _localizedValues[locale.languageCode]?['visitsToday'] ?? 'Visits Today';
  String get errorLoadingDashboard =>
      _localizedValues[locale.languageCode]?['errorLoadingDashboard'] ??
      'Error loading dashboard';
  String get unknownError =>
      _localizedValues[locale.languageCode]?['unknownError'] ?? 'Unknown error';
  String get targetProgress =>
      _localizedValues[locale.languageCode]?['targetProgress'] ??
      'Target Progress';
  String get deals =>
      _localizedValues[locale.languageCode]?['deals'] ?? 'Deals';
  String get viewAll =>
      _localizedValues[locale.languageCode]?['viewAll'] ?? 'View All';
  String get noDataAvailable =>
      _localizedValues[locale.languageCode]?['noDataAvailable'] ??
      'No data available';
  String get pullDownToRefresh =>
      _localizedValues[locale.languageCode]?['pullDownToRefresh'] ??
      'Pull down to refresh';
  String get tomorrow =>
      _localizedValues[locale.languageCode]?['tomorrow'] ?? 'Tomorrow';
  String get yesterday =>
      _localizedValues[locale.languageCode]?['yesterday'] ?? 'Yesterday';
  String get due => _localizedValues[locale.languageCode]?['due'] ?? 'Due';
  String get salesOverview =>
      _localizedValues[locale.languageCode]?['salesOverview'] ??
      'Sales Overview';

  // Accounts & Contacts
  String get searchAccounts =>
      _localizedValues[locale.languageCode]?['searchAccounts'] ??
      'Search accounts...';
  String get searchContacts =>
      _localizedValues[locale.languageCode]?['searchContacts'] ??
      'Search contacts...';
  String get noAccountsFound =>
      _localizedValues[locale.languageCode]?['noAccountsFound'] ??
      'No accounts found';
  String get noContactsFound =>
      _localizedValues[locale.languageCode]?['noContactsFound'] ??
      'No contacts found';
  String get createAccount =>
      _localizedValues[locale.languageCode]?['createAccount'] ??
      'Create Account';
  String get createContact =>
      _localizedValues[locale.languageCode]?['createContact'] ??
      'Create Contact';
  String get accountDetails =>
      _localizedValues[locale.languageCode]?['accountDetails'] ??
      'Account Details';
  String get contactDetails =>
      _localizedValues[locale.languageCode]?['contactDetails'] ??
      'Contact Details';
  String get name => _localizedValues[locale.languageCode]?['name'] ?? 'Name';
  String get category =>
      _localizedValues[locale.languageCode]?['category'] ?? 'Category';
  String get address =>
      _localizedValues[locale.languageCode]?['address'] ?? 'Address';
  String get city => _localizedValues[locale.languageCode]?['city'] ?? 'City';
  String get province =>
      _localizedValues[locale.languageCode]?['province'] ?? 'Province';
  String get phone =>
      _localizedValues[locale.languageCode]?['phone'] ?? 'Phone';
  String get status =>
      _localizedValues[locale.languageCode]?['status'] ?? 'Status';
  String get priority =>
      _localizedValues[locale.languageCode]?['priority'] ?? 'Priority';
  String get urgent =>
      _localizedValues[locale.languageCode]?['urgent'] ?? 'Urgent';
  String get high => _localizedValues[locale.languageCode]?['high'] ?? 'High';
  String get medium =>
      _localizedValues[locale.languageCode]?['medium'] ?? 'Medium';
  String get low => _localizedValues[locale.languageCode]?['low'] ?? 'Low';
  String get general =>
      _localizedValues[locale.languageCode]?['general'] ?? 'General';
  String get call => _localizedValues[locale.languageCode]?['call'] ?? 'Call';
  String get meeting =>
      _localizedValues[locale.languageCode]?['meeting'] ?? 'Meeting';
  String get followUp =>
      _localizedValues[locale.languageCode]?['followUp'] ?? 'Follow Up';
  String get inactive =>
      _localizedValues[locale.languageCode]?['inactive'] ?? 'inactive';
  String get position =>
      _localizedValues[locale.languageCode]?['position'] ?? 'Position';
  String get role => _localizedValues[locale.languageCode]?['role'] ?? 'Role';
  String get notes =>
      _localizedValues[locale.languageCode]?['notes'] ?? 'Notes';
  String get save => _localizedValues[locale.languageCode]?['save'] ?? 'Save';
  String get retry =>
      _localizedValues[locale.languageCode]?['retry'] ?? 'Retry';
  String get viewContacts =>
      _localizedValues[locale.languageCode]?['viewContacts'] ?? 'View Contacts';
  String get selectCategory =>
      _localizedValues[locale.languageCode]?['selectCategory'] ??
      'Select Category';
  String minCharacters(int count) =>
      _localizedValues[locale.languageCode]?['minCharacters']?.replaceAll(
        '{count}',
        count.toString(),
      ) ??
      'Minimum $count characters';
  String maxCharacters(int count) =>
      _localizedValues[locale.languageCode]?['maxCharacters']?.replaceAll(
        '{count}',
        count.toString(),
      ) ??
      'Maximum $count characters';
  String get selectAccount =>
      _localizedValues[locale.languageCode]?['selectAccount'] ??
      'Select Account';
  String get selectRole =>
      _localizedValues[locale.languageCode]?['selectRole'] ?? 'Select Role';
  String get required =>
      _localizedValues[locale.languageCode]?['required'] ?? 'Required';
  String get optional =>
      _localizedValues[locale.languageCode]?['optional'] ?? 'Optional';
  String get edit => _localizedValues[locale.languageCode]?['edit'] ?? 'Edit';
  String get delete =>
      _localizedValues[locale.languageCode]?['delete'] ?? 'Delete';
  String get editAccount =>
      _localizedValues[locale.languageCode]?['editAccount'] ?? 'Edit Account';
  String get deleteAccount =>
      _localizedValues[locale.languageCode]?['deleteAccount'] ??
      'Delete Account';
  String get deleteConfirmation =>
      _localizedValues[locale.languageCode]?['deleteConfirmation'] ??
      'Are you sure you want to delete this account?';
  String get accountDeleted =>
      _localizedValues[locale.languageCode]?['accountDeleted'] ??
      'Account deleted successfully';
  String get accountUpdated =>
      _localizedValues[locale.languageCode]?['accountUpdated'] ??
      'Account updated successfully';
  String get accountCreatedSuccessfully =>
      _localizedValues[locale.languageCode]?['accountCreatedSuccessfully'] ??
      'Account created successfully';
  String get editContact =>
      _localizedValues[locale.languageCode]?['editContact'] ?? 'Edit Contact';
  String get deleteContact =>
      _localizedValues[locale.languageCode]?['deleteContact'] ??
      'Delete Contact';
  String get deleteContactConfirmation =>
      _localizedValues[locale.languageCode]?['deleteContactConfirmation'] ??
      'Are you sure you want to delete this contact?';
  String get contactDeleted =>
      _localizedValues[locale.languageCode]?['contactDeleted'] ??
      'Contact deleted successfully';
  String get contactUpdatedSuccessfully =>
      _localizedValues[locale.languageCode]?['contactUpdatedSuccessfully'] ??
      'Contact updated successfully';
  String get contactCreatedSuccessfully =>
      _localizedValues[locale.languageCode]?['contactCreatedSuccessfully'] ??
      'Contact created successfully';
  String get confirm =>
      _localizedValues[locale.languageCode]?['confirm'] ?? 'Confirm';

  // Route Optimization
  String get routeOptimization =>
      _localizedValues[locale.languageCode]?['routeOptimization'] ??
      'Route Optimization';
  String get createRoute =>
      _localizedValues[locale.languageCode]?['createRoute'] ?? 'Create Route';
  String get routeDetails =>
      _localizedValues[locale.languageCode]?['routeDetails'] ?? 'Route Details';
  String get routeName =>
      _localizedValues[locale.languageCode]?['routeName'] ?? 'Route Name';
  String get routeNameOptional =>
      _localizedValues[locale.languageCode]?['routeNameOptional'] ??
      'Route Name (Optional)';
  String get enterRouteName =>
      _localizedValues[locale.languageCode]?['enterRouteName'] ??
      'Enter route name';
  String get startLocation =>
      _localizedValues[locale.languageCode]?['startLocation'] ??
      'Start Location';
  String get notSet =>
      _localizedValues[locale.languageCode]?['notSet'] ?? 'Not set';
  String get useCurrentLocation =>
      _localizedValues[locale.languageCode]?['useCurrentLocation'] ??
      'Use current location';
  String get waypoints =>
      _localizedValues[locale.languageCode]?['waypoints'] ?? 'Waypoints';
  String get add => _localizedValues[locale.languageCode]?['add'] ?? 'Add';
  String get noWaypointsAdded =>
      _localizedValues[locale.languageCode]?['noWaypointsAdded'] ??
      'No waypoints added';
  String get optimize =>
      _localizedValues[locale.languageCode]?['optimize'] ?? 'Optimize';
  String get noRoutesFound =>
      _localizedValues[locale.languageCode]?['noRoutesFound'] ??
      'No routes found';
  String get createFirstOptimizedRoute =>
      _localizedValues[locale.languageCode]?['createFirstOptimizedRoute'] ??
      'Create your first optimized route';
  String get deleteRoute =>
      _localizedValues[locale.languageCode]?['deleteRoute'] ?? 'Delete Route';
  String get deleteRouteConfirmation =>
      _localizedValues[locale.languageCode]?['deleteRouteConfirmation'] ??
      'Are you sure you want to delete this route? This action cannot be undone.';
  String get routeDeleted =>
      _localizedValues[locale.languageCode]?['routeDeleted'] ??
      'Route deleted successfully';
  String get failedToDeleteRoute =>
      _localizedValues[locale.languageCode]?['failedToDeleteRoute'] ??
      'Failed to delete route';
  String get failedToOptimizeRoute =>
      _localizedValues[locale.languageCode]?['failedToOptimizeRoute'] ??
      'Failed to optimize route';
  String get routeOptimizedSuccessfully =>
      _localizedValues[locale.languageCode]?['routeOptimizedSuccessfully'] ??
      'Route optimized successfully';
  String get unnamedRoute =>
      _localizedValues[locale.languageCode]?['unnamedRoute'] ?? 'Unnamed Route';
  String get distance =>
      _localizedValues[locale.languageCode]?['distance'] ?? 'Distance';
  String get duration =>
      _localizedValues[locale.languageCode]?['duration'] ?? 'Duration';
  String get stops =>
      _localizedValues[locale.languageCode]?['stops'] ?? 'Stops';
  String get routeSteps =>
      _localizedValues[locale.languageCode]?['routeSteps'] ?? 'Route Steps';
  String get continueRoute =>
      _localizedValues[locale.languageCode]?['continueRoute'] ?? 'Continue';
  String get routeMenu =>
      _localizedValues[locale.languageCode]?['routeMenu'] ?? 'Route';
  String get selectWaypoint =>
      _localizedValues[locale.languageCode]?['selectWaypoint'] ??
      'Select Waypoint';
  String get selectFromAccounts =>
      _localizedValues[locale.languageCode]?['selectFromAccounts'] ??
      'Select from Accounts';
  String get selectFromVisitReports =>
      _localizedValues[locale.languageCode]?['selectFromVisitReports'] ??
      'Select from Visit Reports';
  String get searchAccountsOrContacts =>
      _localizedValues[locale.languageCode]?['searchAccountsOrContacts'] ??
      'Search accounts or contacts...';
  String get noAccountsFoundForWaypoint =>
      _localizedValues[locale.languageCode]?['noAccountsFoundForWaypoint'] ??
      'No accounts found';
  String get noVisitReportsFoundForWaypoint =>
      _localizedValues[locale
          .languageCode]?['noVisitReportsFoundForWaypoint'] ??
      'No visit reports found';
  String get destination =>
      _localizedValues[locale.languageCode]?['destination'] ?? 'Destination';

  // Leads
  String get leads =>
      _localizedValues[locale.languageCode]?['leads'] ?? 'Leads';

  String get pipeline =>
      _localizedValues[locale.languageCode]?['pipeline'] ?? 'Pipeline';

  String get createLead =>
      _localizedValues[locale.languageCode]?['createLead'] ?? 'Create Lead';
  String get editLead =>
      _localizedValues[locale.languageCode]?['editLead'] ?? 'Edit Lead';
  String get leadDetails =>
      _localizedValues[locale.languageCode]?['leadDetails'] ?? 'Lead Details';
  String get searchLeads =>
      _localizedValues[locale.languageCode]?['searchLeads'] ??
      'Search leads...';
  String get noLeadsFound =>
      _localizedValues[locale.languageCode]?['noLeadsFound'] ??
      'No leads found';
  String get deleteLead =>
      _localizedValues[locale.languageCode]?['deleteLead'] ?? 'Delete Lead';
  String get deleteLeadConfirmation =>
      _localizedValues[locale.languageCode]?['deleteLeadConfirmation'] ??
      'Are you sure you want to delete this lead?';
  String get leadDeletedSuccessfully =>
      _localizedValues[locale.languageCode]?['leadDeletedSuccessfully'] ??
      'Lead deleted successfully';
  String get failedToDeleteLead =>
      _localizedValues[locale.languageCode]?['failedToDeleteLead'] ??
      'Failed to delete lead';
  String get leadCreatedSuccessfully =>
      _localizedValues[locale.languageCode]?['leadCreatedSuccessfully'] ??
      'Lead created successfully';
  String get leadUpdatedSuccessfully =>
      _localizedValues[locale.languageCode]?['leadUpdatedSuccessfully'] ??
      'Lead updated successfully';
  String get failedToCreateLead =>
      _localizedValues[locale.languageCode]?['failedToCreateLead'] ??
      'Failed to create lead';
  String get failedToUpdateLead =>
      _localizedValues[locale.languageCode]?['failedToUpdateLead'] ??
      'Failed to update lead';
  String get convertLead =>
      _localizedValues[locale.languageCode]?['convertLead'] ?? 'Convert Lead';
  String get convertLeadConfirmation =>
      _localizedValues[locale.languageCode]?['convertLeadConfirmation'] ??
      'Are you sure you want to convert this lead?';
  String get leadConvertedSuccessfully =>
      _localizedValues[locale.languageCode]?['leadConvertedSuccessfully'] ??
      'Lead converted successfully';
  String get failedToConvertLead =>
      _localizedValues[locale.languageCode]?['failedToConvertLead'] ??
      'Failed to convert lead';
  String get company =>
      _localizedValues[locale.languageCode]?['company'] ?? 'Company';
  String get industry =>
      _localizedValues[locale.languageCode]?['industry'] ?? 'Industry';
  String get source =>
      _localizedValues[locale.languageCode]?['source'] ?? 'Source';
  String get website =>
      _localizedValues[locale.languageCode]?['website'] ?? 'Website';

  String get probability =>
      _localizedValues[locale.languageCode]?['probability'] ??
      'Probability (%)';

  String get opportunityTitle =>
      _localizedValues[locale.languageCode]?['opportunityTitle'] ??
      'Opportunity Title';

  String get pipelineStage =>
      _localizedValues[locale.languageCode]?['pipelineStage'] ??
      'Pipeline Stage';

  String get dealValue =>
      _localizedValues[locale.languageCode]?['dealValue'] ?? 'Deal Value';

  String get expectedCloseDate =>
      _localizedValues[locale.languageCode]?['expectedCloseDate'] ??
      'Expected Close Date';

  String get convert =>
      _localizedValues[locale.languageCode]?['convert'] ?? 'Convert';

  String get isRequired =>
      _localizedValues[locale.languageCode]?['isRequired'] ?? 'is required';

  String get clear =>
      _localizedValues[locale.languageCode]?['clear'] ?? 'Clear';

  String get apply =>
      _localizedValues[locale.languageCode]?['apply'] ?? 'Apply';

  String get filterBySource =>
      _localizedValues[locale.languageCode]?['filterBySource'] ??
      'Filter by Source';

  String get filterByIndustry =>
      _localizedValues[locale.languageCode]?['filterByIndustry'] ??
      'Filter by Industry';

  String get filterByProvince =>
      _localizedValues[locale.languageCode]?['filterByProvince'] ??
      'Filter by Province';

  static final Map<String, Map<String, String>> _localizedValues = {
    'en': {
      'profile': 'Profile',
      'settings': 'Settings',
      'notifications': 'Notifications',
      'manageNotificationSettings': 'Manage notification settings',
      'language': 'Language',
      'theme': 'Theme',
      'lightTheme': 'Light theme',
      'darkTheme': 'Dark theme',
      'systemTheme': 'Follow system setting',
      'about': 'About',
      'appVersion': 'App Version',
      'privacyPolicy': 'Privacy Policy',
      'viewPrivacyPolicy': 'View privacy policy',
      'termsOfService': 'Terms of Service',
      'viewTermsOfService': 'View terms of service',
      'logout': 'Logout',
      'logoutConfirmation': 'Are you sure you want to logout?',
      'cancel': 'Cancel',
      'selectTheme': 'Select Theme',
      'selectLanguage': 'Select Language',
      'english': 'English',
      'indonesian': 'Indonesian',
      'signInToYourAccount': 'Sign in to your Account',
      'enterEmailPassword': 'Enter your email and password to log in',
      'email': 'Email',
      'password': 'Password',
      'enterPassword': 'Enter your password',
      'rememberMe': 'Remember me',
      'logIn': 'Log In',
      'home': 'Home',
      'accounts': 'Accounts',
      'contacts': 'Contacts',
      'accountsAndContacts': 'Accounts',
      'reports': 'Reports',
      'tasks': 'Tasks',
      'reportsAndTasks': 'Reports',
      'visitsMenu': 'Visits',
      'search': 'Search',
      'searchVisitReports': 'Search visit reports...',
      'searchTasks': 'Search tasks...',
      'noVisitReportsFound': 'No visit reports found',
      'noTasksFound': 'No tasks found',
      'noSearchResults': 'No results found for',
      'searchResultsFor': 'Search results for',
      'tapToCreateVisitReport': 'Tap + to create a new visit report',
      'tapToCreateTask': 'Tap + to create a new task',
      'createVisitReport': 'Create Visit Report',
      'visitReports': 'Visit Reports',
      'visitReportDetails': 'Visit Report Details',
      'visitInformation': 'Visit Information',
      'checkInOutStatus': 'Check-in/out Status',
      'visitDate': 'Visit Date',
      'purpose': 'Purpose',
      'photos': 'Photos',
      'checkInTime': 'Check-in Time',
      'checkOutTime': 'Check-out Time',
      'checkInLocation': 'Check-in Location',
      'checkOutLocation': 'Check-out Location',
      'checkIn': 'Check In',
      'checkOut': 'Check Out',
      'uploadPhoto': 'Upload Photo',
      'notCheckedIn': 'Not checked in',
      'notCheckedOut': 'Not checked out',
      'checkInSuccessful': 'Check-in successful',
      'checkOutSuccessful': 'Check-out successful',
      'photoUploadedSuccessfully': 'Photo uploaded successfully',
      'failedToCheckIn': 'Failed to check in',
      'failedToCheckOut': 'Failed to check out',
      'failedToUploadPhoto': 'Failed to upload photo',
      'selfieRequiredForCheckIn':
          'Selfie picture is required for check-in. Please take a photo to continue.',
      'previewSelfie': 'Preview Selfie',
      'retake': 'Retake',
      'confirmCheckOut': 'Confirm Check-out',
      'checkOutLocationRequired': 'Current location is required for check-out',
      'gettingLocation': 'Getting your current location...',
      'locationAccuracy': 'Location Accuracy',
      'currentLocation': 'Current Location',
      'fakeGPSDetected': 'Fake GPS Detected',
      'fakeGPSDescription':
          'We detected that you are using a Fake GPS application. Check-in and check-out features are disabled when Fake GPS is active.',
      'detectedReason': 'Detected Reason',
      'howToDisableFakeGPS': 'How to Disable Fake GPS',
      'fakeGPSStep1':
          'Close or uninstall any Fake GPS or location spoofing applications on your device',
      'fakeGPSStep2':
          'Restart your device to ensure all location services are reset',
      'fakeGPSStep3': 'Clear app cache and location permissions',
      'fakeGPSStep4': 'Try check-in/check-out again after disabling Fake GPS',
      'fakeGPSImportantNote':
          'Check-in and check-out features require real GPS location for verification. Using Fake GPS is not allowed and will prevent you from using these features.',
      'close': 'Close',
      'selectContact': 'Select Contact',
      'selectContactOptional': 'Select contact (optional)',
      'selectVisitDate': 'Select visit date',
      'enterVisitPurpose': 'Enter visit purpose',
      'enterAdditionalNotes': 'Enter additional notes',
      'pleaseSelectAccount': 'Please select an account',
      'deal': 'Deal',
      'lead': 'Lead',
      'none': 'None',
      'pleaseSelectDeal': 'Please select a deal',
      'pleaseSelectLead': 'Please select a lead',
      'noLeadsAvailable': 'No leads available',
      'purposeRequired': 'Purpose is required',
      'updateVisitReport': 'Update Visit Report',
      'deleteVisitReport': 'Delete Visit Report',
      'submitVisitReport': 'Submit Visit Report',
      'deleteVisitReportConfirmation':
          'Are you sure you want to delete this visit report? This action cannot be undone.',
      'visitReportUpdatedSuccessfully': 'Visit report updated successfully',
      'visitReportDeletedSuccessfully': 'Visit report deleted successfully',
      'visitReportSubmittedSuccessfully': 'Visit report submitted successfully',
      'failedToUpdateVisitReport': 'Failed to update visit report',
      'failedToDeleteVisitReport': 'Failed to delete visit report',
      'failedToSubmitVisitReport': 'Failed to submit visit report',
      'outcome': 'Outcome',
      'nextSteps': 'Next Steps',
      'selectOutcome': 'Select outcome',
      'enterNextSteps': 'Enter next steps (optional)',
      'visitReportCreatedSuccessfully': 'Visit report created successfully',
      'failedToCreateVisitReport': 'Failed to create visit report',
      'all': 'All',
      'filterByStatus': 'Filter by Status',
      // Route Optimization
      'routeOptimization': 'Route Optimization',
      'createRoute': 'Create Route',
      'routeDetails': 'Route Details',
      'routeName': 'Route Name',
      'routeNameOptional': 'Route Name (Optional)',
      'enterRouteName': 'Enter route name',
      'startLocation': 'Start Location',
      'notSet': 'Not set',
      'useCurrentLocation': 'Use current location',
      'waypoints': 'Waypoints',
      'add': 'Add',
      'noWaypointsAdded': 'No waypoints added',
      'optimize': 'Optimize',
      'noRoutesFound': 'No routes found',
      'createFirstOptimizedRoute': 'Create your first optimized route',
      'deleteRoute': 'Delete Route',
      'deleteRouteConfirmation':
          'Are you sure you want to delete this route? This action cannot be undone.',
      'routeDeleted': 'Route deleted successfully',
      'failedToDeleteRoute': 'Failed to delete route',
      'failedToOptimizeRoute': 'Failed to optimize route',
      'routeOptimizedSuccessfully': 'Route optimized successfully',
      'unnamedRoute': 'Unnamed Route',
      'distance': 'Distance',
      'duration': 'Duration',
      'stops': 'Stops',
      'routeSteps': 'Route Steps',
      'continueRoute': 'Continue',
      'routeMenu': 'Route',
      'selectWaypoint': 'Select Waypoint',
      'selectFromAccounts': 'Select from Accounts',
      'selectFromVisitReports': 'Select from Visit Reports',
      'searchAccountsOrContacts': 'Search accounts or contacts...',
      'noAccountsFoundForWaypoint': 'No accounts found',
      'noVisitReportsFoundForWaypoint': 'No visit reports found',
      'destination': 'Destination',
      'filterByPriority': 'Filter by Priority',
      'filterByType': 'Filter by Type',
      'filterByDueDate': 'Filter by Due Date',
      'dueDateFrom': 'Due Date From',
      'dueDateTo': 'Due Date To',
      'selectDueDateFrom': 'Select due date from',
      'selectDueDateTo': 'Select due date to',
      'clearFilters': 'Clear filters',
      'clearSearch': 'Clear Search',
      'taskDetails': 'Task Details',
      'taskInformation': 'Task Information',
      'relatedInformation': 'Related Information',
      'reminders': 'Reminders',
      'completeTask': 'Complete Task',
      'markInProgress': 'Mark In Progress',
      'markInProgressConfirmation':
          'Are you sure you want to mark this task as in progress?',
      'taskMarkedInProgress': 'Task marked as in progress',
      'failedToMarkInProgress': 'Failed to mark task as in progress',
      'addReminder': 'Add Reminder',
      'completeTaskConfirmation':
          'Are you sure you want to mark this task as completed?',
      'taskCompletedSuccessfully': 'Task completed successfully',
      'failedToCompleteTask': 'Failed to complete task',
      'deleteTask': 'Delete Task',
      'deleteTaskConfirmation':
          'Are you sure you want to delete this task? This action cannot be undone.',
      'taskDeletedSuccessfully': 'Task deleted successfully',
      'failedToDeleteTask': 'Failed to delete task',
      'reminderCreatedSuccessfully': 'Reminder created successfully',
      'failedToCreateReminder': 'Failed to create reminder',
      'selectReminderDate': 'Select reminder date and time',
      'reminderMessage': 'Message (optional)',
      'enterReminderMessage': 'Enter reminder message',
      'title': 'Title',
      'description': 'Description',
      'type': 'Type',
      'dueDate': 'Due Date',
      'selectDueDate': 'Select due date',
      'createTask': 'Create Task',
      'editTask': 'Edit Task',
      'updateTask': 'Update Task',
      'enterTaskTitle': 'Enter task title',
      'enterTaskDescription': 'Enter task description',
      'taskCreatedSuccessfully': 'Task created successfully',
      'taskUpdatedSuccessfully': 'Task updated successfully',
      'failedToCreateTask': 'Failed to create task',
      'failedToUpdateTask': 'Failed to update task',
      'titleIsRequired': 'Title is required',
      'sent': 'Sent',
      'noNotificationsFound': 'No notifications found',
      'markAllAsRead': 'Mark All as Read',
      'markAsRead': 'Mark as Read',
      'filter': 'Filter',
      'unread': 'Unread',
      'read': 'Read',
      'filterNotifications': 'Filter Notifications',
      'notificationMarkedAsRead': 'Notification marked as read',
      'failedToMarkAsRead': 'Failed to mark as read',
      'allNotificationsMarkedAsRead': 'All notifications marked as read',
      'failedToMarkAllAsRead': 'Failed to mark all as read',
      'deleteNotification': 'Delete Notification',
      'deleteNotificationConfirmation':
          'Are you sure you want to delete this notification?',
      'notificationDeleted': 'Notification deleted successfully',
      'failedToDeleteNotification': 'Failed to delete notification',
      'dashboard': 'Dashboard',
      'welcomeBack': 'Welcome back,',
      'today': 'Today',
      'thisWeek': 'This Week',
      'thisMonth': 'This Month',
      'thisYear': 'This Year',
      'totalVisits': 'Total Visits',
      'totalAccounts': 'Total Accounts',
      'totalActivities': 'Total Activities',
      'revenue': 'Revenue',
      'completed': 'completed',
      'planned': 'Planned',
      'cancelled': 'Cancelled',
      'seeAll': 'See All',
      'noVisitsFound': 'No visits found',
      'remaining': 'remaining',
      'tersisa': 'tersisa',
      'active': 'active',
      'visits': 'visits',
      'calls': 'calls',
      'fromWonDeals': 'From won deals',
      'salesTarget': 'Sales Target',
      'totalDeals': 'Total Deals',
      'leadsBySource': 'Leads by Source',
      'upcomingTasks': 'Upcoming Tasks',
      'pipelineSummary': 'Pipeline Summary',
      'salesPipeline': 'Sales Pipeline',
      'recentActivities': 'Recent Activities',
      'target': 'Target',
      'achieved': 'Achieved',
      'open': 'Open',
      'won': 'Won',
      'lost': 'Lost',
      'totalValue': 'Total Value',
      'totalLeads': 'total leads',
      'noLeadsForPeriod': 'No leads for this period',
      'noUpcomingTasks': 'No upcoming tasks',
      'targetProgressDescription': '{progress}% of target achieved',
      'totalAccountsDescription': '{active} active, {inactive} inactive',
      'totalDealsDescription': '{open} open, {won} won',
      'totalRevenueDescription': 'Based on won deals in this period',
      'topAccounts': 'Top Accounts',
      'topSalesReps': 'Top Sales Reps',
      'visitStatistics': 'Visit Statistics',
      'activityTrends': 'Activity Trends',
      'noTopAccounts': 'No accounts found',
      'noTopSalesReps': 'No sales reps found',
      'topAccountsVisits': '{count} visits',
      'topAccountsActivities': '{count} activities',
      'topSalesRepsVisits': '{count} visits',
      'topSalesRepsAccounts': '{count} accounts',
      'emails': 'emails',
      'pending': 'Pending',
      'confirmed': 'Confirmed',
      'approved': 'approved',
      'draft': 'Draft',
      'submitted': 'Submitted',
      'rejected': 'Rejected',
      'total': 'Total',
      'visitsToday': 'Visits Today',
      'errorLoadingDashboard': 'Error loading dashboard',
      'unknownError': 'Unknown error',
      'targetProgress': 'Target Progress',
      'deals': 'Deals',
      'viewAll': 'View All',
      'noDataAvailable': 'No data available',
      'pullDownToRefresh': 'Pull down to refresh',
      'tomorrow': 'Tomorrow',
      'yesterday': 'Yesterday',
      'due': 'Due',
      'salesOverview': 'Sales Overview',
      'searchAccounts': 'Search accounts...',
      'searchContacts': 'Search contacts...',
      'noAccountsFound': 'No accounts found',
      'noContactsFound': 'No contacts found',
      'createAccount': 'Create Account',
      'createContact': 'Create Contact',
      'accountDetails': 'Account Details',
      'contactDetails': 'Contact Details',
      'name': 'Name',
      'category': 'Category',
      'address': 'Address',
      'city': 'City',
      'province': 'Province',
      'phone': 'Phone',
      'status': 'Status',
      'priority': 'Priority',
      'urgent': 'Urgent',
      'high': 'High',
      'medium': 'Medium',
      'low': 'Low',
      'general': 'General',
      'call': 'Call',
      'meeting': 'Meeting',
      'followUp': 'Follow Up',
      'inactive': 'inactive',
      'position': 'Position',
      'role': 'Role',
      'notes': 'Notes',
      'scheduled': 'Scheduled',
      'save': 'Save',
      'retry': 'Retry',
      'viewContacts': 'View Contacts',
      'selectCategory': 'Select Category',
      'selectAccount': 'Select Account',
      'selectRole': 'Select Role',
      'required': 'Required',
      'optional': 'Optional',
      'edit': 'Edit',
      'delete': 'Delete',
      'editAccount': 'Edit Account',
      'deleteAccount': 'Delete Account',
      'deleteConfirmation': 'Are you sure you want to delete this account?',
      'accountDeleted': 'Account deleted successfully',
      'accountUpdated': 'Account updated successfully',
      'accountCreatedSuccessfully': 'Account created successfully',
      'editContact': 'Edit Contact',
      'deleteContact': 'Delete Contact',
      'deleteContactConfirmation':
          'Are you sure you want to delete this contact?',
      'contactDeleted': 'Contact deleted successfully',
      'contactUpdatedSuccessfully': 'Contact updated successfully',
      'contactCreatedSuccessfully': 'Contact created successfully',
      'confirm': 'Confirm',
      // Leads
      'leads': 'Leads',
      'pipeline': 'Pipeline',
      'createLead': 'Create Lead',
      'editLead': 'Edit Lead',
      'leadDetails': 'Lead Details',
      'searchLeads': 'Search leads...',
      'noLeadsFound': 'No leads found',
      'deleteLead': 'Delete Lead',
      'deleteLeadConfirmation': 'Are you sure you want to delete this lead?',
      'leadDeletedSuccessfully': 'Lead deleted successfully',
      'failedToDeleteLead': 'Failed to delete lead',
      'leadCreatedSuccessfully': 'Lead created successfully',
      'leadUpdatedSuccessfully': 'Lead updated successfully',
      'firstName': 'First Name',
      'lastName': 'Last Name',
      'failedToCreateLead': 'Failed to create lead',
      'failedToUpdateLead': 'Failed to update lead',
      'convertLead': 'Convert Lead',
      'convertLeadConfirmation': 'Are you sure you want to convert this lead?',
      'leadConvertedSuccessfully': 'Lead converted successfully',
      'failedToConvertLead': 'Failed to convert lead',
      'company': 'Company',
      'industry': 'Industry',
      'source': 'Source',
      'website': 'Website',
      'convert': 'Convert',
      'opportunityTitle': 'Opportunity Title',
      'pipelineStage': 'Pipeline Stage',
      'dealValue': 'Deal Value',
      'expectedCloseDate': 'Expected Close Date',
      'probability': 'Probability (%)',
      'isRequired': 'is required',
      'clear': 'Clear',
      'apply': 'Apply',
      'filterBySource': 'Filter by Source',
      'filterByIndustry': 'Filter by Industry',
      'filterByProvince': 'Filter by Province',
      'searchSchedules': 'Search schedules...',
      'noSchedulesFound': 'No schedules found',
      'tapToCreateSchedule': 'Tap + to create a new schedule',
      'schedules': 'Schedules',
      'scheduleDetails': 'Schedule Details',
      'scheduleInformation': 'Schedule Information',
      'scheduledAt': 'Scheduled At',
      'reminder': 'Reminder',
      'minutesBefore': 'minutes before',
      'linkedTask': 'Linked Task',
      'selectTask': 'Select Task',
      'scheduledAtRequired': 'Scheduled time is required',
      'scheduledFrom': 'Scheduled From',
      'scheduledTo': 'Scheduled To',
      'selectDate': 'Select Date',
      'clearDate': 'Clear Date',
      'selectScheduledFrom': 'Select scheduled from',
      'selectScheduledTo': 'Select scheduled to',
      'noData': 'No data',
      'noDescription': 'No description',
      'location': 'Location',
      'noLocation': 'No location',
      'relatedVisitReport': 'Related Visit Report',
      'viewVisitReport': 'View Visit Report',
      'activityType': 'Activity Type',
      'createSchedule': 'Create Schedule',
      'editSchedule': 'Edit Schedule',
      'updateSchedule': 'Update Schedule',
      'deleteSchedule': 'Delete Schedule',
      'deleteScheduleConfirmation':
          'Are you sure you want to delete this schedule? This action cannot be undone.',
      'scheduleCreatedSuccessfully': 'Schedule created successfully',
      'scheduleUpdatedSuccessfully': 'Schedule updated successfully',
      'scheduleDeletedSuccessfully': 'Schedule deleted successfully',
      'failedToCreateSchedule': 'Failed to create schedule',
      'failedToUpdateSchedule': 'Failed to update schedule',
      'failedToDeleteSchedule': 'Failed to delete schedule',
      'reminderMinutesBefore': 'Reminder (minutes before)',
      'invalidReminderMinutes': 'Invalid reminder minutes (0-10080)',
      'enterTitle': 'Enter title',
      'minCharacters': 'Minimum {count} characters',
      'maxCharacters': 'Maximum {count} characters',
    },
    'id': {
      'profile': 'Profil',
      'settings': 'Pengaturan',
      'notifications': 'Notifikasi',
      'manageNotificationSettings': 'Kelola pengaturan notifikasi',
      'language': 'Bahasa',
      'theme': 'Tema',
      'lightTheme': 'Tema terang',
      'darkTheme': 'Tema gelap',
      'systemTheme': 'Ikuti pengaturan sistem',
      'about': 'Tentang',
      'appVersion': 'Versi Aplikasi',
      'privacyPolicy': 'Kebijakan Privasi',
      'viewPrivacyPolicy': 'Lihat kebijakan privasi',
      'termsOfService': 'Syarat Layanan',
      'viewTermsOfService': 'Lihat syarat layanan',
      'logout': 'Keluar',
      'logoutConfirmation': 'Apakah Anda yakin ingin keluar?',
      'cancel': 'Batal',
      'selectTheme': 'Pilih Tema',
      'selectLanguage': 'Pilih Bahasa',
      'english': 'Bahasa Inggris',
      'indonesian': 'Bahasa Indonesia',
      'signInToYourAccount': 'Masuk ke Akun Anda',
      'enterEmailPassword': 'Masukkan email dan kata sandi Anda untuk masuk',
      'email': 'Email',
      'password': 'Kata Sandi',
      'enterPassword': 'Masukkan kata sandi Anda',
      'rememberMe': 'Ingat saya',
      'logIn': 'Masuk',
      'home': 'Beranda',
      'accounts': 'Akun',
      'contacts': 'Kontak',
      'accountsAndContacts': 'Akun',
      'reports': 'Laporan',
      'tasks': 'Tugas',
      'reportsAndTasks': 'Laporan',
      'visitsMenu': 'Visits', // Tetap 'Visits' untuk bahasa Indonesia
      'clearFilters': 'Hapus filter',
      'filterByStatus': 'Filter berdasarkan Status',
      'filterByPriority': 'Filter berdasarkan Prioritas',
      'filterByType': 'Filter berdasarkan Tipe',
      'filterByDueDate': 'Filter berdasarkan Tanggal Jatuh Tempo',
      'dueDateFrom': 'Tanggal Jatuh Tempo Dari',
      'dueDateTo': 'Tanggal Jatuh Tempo Sampai',
      'selectDueDateFrom': 'Pilih tanggal jatuh tempo dari',
      'selectDueDateTo': 'Pilih tanggal jatuh tempo sampai',
      'taskDetails': 'Detail Tugas',
      'taskInformation': 'Informasi Tugas',
      'relatedInformation': 'Informasi Terkait',
      'reminders': 'Pengingat',
      'completeTask': 'Selesaikan Tugas',
      'markInProgress': 'Tandai Sedang Berjalan',
      'markInProgressConfirmation':
          'Apakah Anda yakin ingin menandai tugas ini sebagai sedang berjalan?',
      'taskMarkedInProgress': 'Tugas ditandai sebagai sedang berjalan',
      'failedToMarkInProgress': 'Gagal menandai tugas sebagai sedang berjalan',
      'addReminder': 'Tambah Pengingat',
      'completeTaskConfirmation':
          'Apakah Anda yakin ingin menandai tugas ini sebagai selesai?',
      'taskCompletedSuccessfully': 'Tugas berhasil diselesaikan',
      'failedToCompleteTask': 'Gagal menyelesaikan tugas',
      'deleteTask': 'Hapus Tugas',
      'deleteTaskConfirmation':
          'Apakah Anda yakin ingin menghapus tugas ini? Tindakan ini tidak dapat dibatalkan.',
      'taskDeletedSuccessfully': 'Tugas berhasil dihapus',
      'failedToDeleteTask': 'Gagal menghapus tugas',
      'reminderCreatedSuccessfully': 'Pengingat berhasil dibuat',
      'failedToCreateReminder': 'Gagal membuat pengingat',
      'selectReminderDate': 'Pilih tanggal dan waktu pengingat',
      'reminderMessage': 'Pesan (opsional)',
      'enterReminderMessage': 'Masukkan pesan pengingat',
      'title': 'Judul',
      'description': 'Deskripsi',
      'type': 'Tipe',
      'dueDate': 'Tanggal Jatuh Tempo',
      'selectDueDate': 'Pilih tanggal jatuh tempo',
      'createTask': 'Buat Tugas',
      'editTask': 'Edit Tugas',
      'updateTask': 'Perbarui Tugas',
      'enterTaskTitle': 'Masukkan judul tugas',
      'enterTaskDescription': 'Masukkan deskripsi tugas',
      'taskCreatedSuccessfully': 'Tugas berhasil dibuat',
      'taskUpdatedSuccessfully': 'Tugas berhasil diperbarui',
      'failedToCreateTask': 'Gagal membuat tugas',
      'failedToUpdateTask': 'Gagal memperbarui tugas',
      'titleIsRequired': 'Judul wajib diisi',
      'sent': 'Terkirim',
      'markAllAsRead': 'Tandai Semua sebagai Dibaca',
      'markAsRead': 'Tandai sebagai Dibaca',
      'filter': 'Filter',
      'unread': 'Belum Dibaca',
      'read': 'Dibaca',
      'filterNotifications': 'Filter Notifikasi',
      'notificationMarkedAsRead': 'Notifikasi ditandai sebagai dibaca',
      'failedToMarkAsRead': 'Gagal menandai sebagai dibaca',
      'allNotificationsMarkedAsRead':
          'Semua notifikasi ditandai sebagai dibaca',
      'failedToMarkAllAsRead': 'Gagal menandai semua sebagai dibaca',
      'deleteNotification': 'Hapus Notifikasi',
      'deleteNotificationConfirmation':
          'Apakah Anda yakin ingin menghapus notifikasi ini?',
      'notificationDeleted': 'Notifikasi berhasil dihapus',
      'failedToDeleteNotification': 'Gagal menghapus notifikasi',
      'dashboard': 'Dashboard',
      'welcomeBack': 'Selamat datang kembali,',
      'today': 'Hari Ini',
      'thisWeek': 'Minggu Ini',
      'thisMonth': 'Bulan Ini',
      'thisYear': 'Tahun Ini',
      'totalVisits': 'Total Kunjungan',
      'totalAccounts': 'Total Akun',
      'totalActivities': 'Total Aktivitas',
      'revenue': 'Pendapatan',
      'completed': 'selesai',
      'planned': 'Terencana',
      'cancelled': 'Dibatalkan',
      'seeAll': 'Lihat Semua',
      'noVisitsFound': 'Tidak ada kunjungan ditemukan',
      'remaining': 'remaining',
      'tersisa': 'tersisa',
      'active': 'aktif',
      'visits': 'visits',
      'calls': 'panggilan',
      'fromWonDeals': 'Dari deal yang menang',
      'salesTarget': 'Target Penjualan',
      'totalDeals': 'Total Deal',
      'leadsBySource': 'Leads berdasarkan Sumber',
      'upcomingTasks': 'Tugas Mendatang',
      'pipelineSummary': 'Ringkasan Pipeline',
      'salesPipeline': 'Pipeline Penjualan',
      'recentActivities': 'Aktivitas Terkini',
      'target': 'Target',
      'achieved': 'Tercapai',
      'open': 'Terbuka',
      'won': 'Menang',
      'lost': 'Kalah',
      'totalValue': 'Total Nilai',
      'totalLeads': 'total leads',
      'noLeadsForPeriod': 'Tidak ada leads di periode ini',
      'noUpcomingTasks': 'Tidak ada tugas mendatang',
      'targetProgressDescription': '{progress}% dari target tercapai',
      'totalAccountsDescription': '{active} aktif, {inactive} tidak aktif',
      'totalDealsDescription': '{open} terbuka, {won} menang',
      'totalRevenueDescription': 'Berdasarkan deal menang di periode ini',
      'topAccounts': 'Akun Teratas',
      'topSalesReps': 'Sales Teratas',
      'visitStatistics': 'Statistik Kunjungan',
      'activityTrends': 'Tren Aktivitas',
      'noTopAccounts': 'Belum ada akun',
      'noTopSalesReps': 'Belum ada data sales',
      'topAccountsVisits': '{count} kunjungan',
      'topAccountsActivities': '{count} aktivitas',
      'topSalesRepsVisits': '{count} kunjungan',
      'topSalesRepsAccounts': '{count} akun',
      'emails': 'email',
      'approved': 'disetujui',
      'draft': 'Draft',
      'submitted': 'Dikirim',
      'rejected': 'Ditolak',
      'total': 'Total',
      'visitsToday': 'Kunjungan Hari Ini',
      'errorLoadingDashboard': 'Gagal memuat dashboard',
      'unknownError': 'Error tidak diketahui',
      'targetProgress': 'Progress Target',
      'deals': 'Deal',
      'viewAll': 'Lihat Semua',
      'noDataAvailable': 'Tidak ada data tersedia',
      'pullDownToRefresh': 'Tarik ke bawah untuk refresh',
      'tomorrow': 'Besok',
      'yesterday': 'Kemarin',
      'due': 'Jatuh Tempo',
      'salesOverview': 'Ringkasan Penjualan',
      'searchAccounts': 'Cari akun...',
      'searchContacts': 'Cari kontak...',
      'noAccountsFound': 'Tidak ada akun ditemukan',
      'noContactsFound': 'Tidak ada kontak ditemukan',
      'createAccount': 'Buat Akun',
      'createContact': 'Buat Kontak',
      'accountDetails': 'Detail Akun',
      'contactDetails': 'Detail Kontak',
      'name': 'Nama',
      'category': 'Kategori',
      'address': 'Alamat',
      'city': 'Kota',
      'province': 'Provinsi',
      'phone': 'Telepon',
      'status': 'Status',
      'priority': 'Prioritas',
      'urgent': 'Mendesak',
      'high': 'Tinggi',
      'medium': 'Sedang',
      'low': 'Rendah',
      'general': 'Umum',
      'call': 'Panggilan',
      'meeting': 'Rapat',
      'followUp': 'Tindak Lanjut',
      'inactive': 'tidak aktif',
      'position': 'Posisi',
      'role': 'Peran',
      'notes': 'Catatan',
      'save': 'Simpan',
      'retry': 'Coba Lagi',
      'viewContacts': 'Lihat Kontak',
      'selectCategory': 'Pilih Kategori',
      'selectAccount': 'Pilih Akun',
      'selectRole': 'Pilih Peran',
      'required': 'Wajib',
      'optional': 'Opsional',
      'edit': 'Edit',
      'delete': 'Hapus',
      'editAccount': 'Edit Akun',
      'deleteAccount': 'Hapus Akun',
      'deleteConfirmation': 'Apakah Anda yakin ingin menghapus akun ini?',
      'accountDeleted': 'Akun berhasil dihapus',
      'accountUpdated': 'Akun berhasil diperbarui',
      'accountCreatedSuccessfully': 'Akun berhasil dibuat',
      'editContact': 'Edit Kontak',
      'deleteContact': 'Hapus Kontak',
      'deleteContactConfirmation':
          'Apakah Anda yakin ingin menghapus kontak ini?',
      'contactDeleted': 'Kontak berhasil dihapus',
      'contactUpdatedSuccessfully': 'Kontak berhasil diperbarui',
      'contactCreatedSuccessfully': 'Kontak berhasil dibuat',
      'confirm': 'Konfirmasi',
      'search': 'Cari',
      'searchVisitReports': 'Cari laporan kunjungan...',
      'searchTasks': 'Cari tugas...',
      'noVisitReportsFound': 'Tidak ada laporan kunjungan ditemukan',
      'noTasksFound': 'Tidak ada tugas ditemukan',
      'noSearchResults': 'Tidak ada hasil ditemukan untuk',
      'searchResultsFor': 'Hasil pencarian untuk',
      'tapToCreateVisitReport': 'Ketuk + untuk membuat laporan kunjungan baru',
      'tapToCreateTask': 'Ketuk + untuk membuat tugas baru',
      'createVisitReport': 'Buat Laporan Kunjungan',
      'visitReports': 'Visit Reports',
      'visitReportDetails': 'Detail Laporan Kunjungan',
      'visitInformation': 'Informasi Kunjungan',
      'checkInOutStatus': 'Status Check-in/out',
      'visitDate': 'Tanggal Kunjungan',
      'purpose': 'Tujuan',
      'photos': 'Foto',
      'checkInTime': 'Waktu Check-in',
      'checkOutTime': 'Waktu Check-out',
      'checkInLocation': 'Lokasi Check-in',
      'checkOutLocation': 'Lokasi Check-out',
      'checkIn': 'Check In',
      'checkOut': 'Check Out',
      'uploadPhoto': 'Unggah Foto',
      'notCheckedIn': 'Belum check in',
      'notCheckedOut': 'Belum check out',
      'checkInSuccessful': 'Check-in berhasil',
      'checkOutSuccessful': 'Check-out berhasil',
      'photoUploadedSuccessfully': 'Foto berhasil diunggah',
      'failedToCheckIn': 'Gagal check in',
      'failedToCheckOut': 'Gagal check out',
      'failedToUploadPhoto': 'Gagal mengunggah foto',
      'selfieRequiredForCheckIn':
          'Foto selfie wajib untuk check-in. Silakan ambil foto untuk melanjutkan.',
      'previewSelfie': 'Pratinjau Selfie',
      'confirmCheckOut': 'Konfirmasi Check-out',
      'checkOutLocationRequired': 'Lokasi saat ini diperlukan untuk check-out',
      'gettingLocation': 'Mendapatkan lokasi Anda...',
      'locationAccuracy': 'Akurasi Lokasi',
      'currentLocation': 'Lokasi Saat Ini',
      'fakeGPSDetected': 'Fake GPS Terdeteksi',
      'fakeGPSDescription':
          'Kami mendeteksi bahwa Anda menggunakan aplikasi Fake GPS. Fitur check-in dan check-out dinonaktifkan saat Fake GPS aktif.',
      'detectedReason': 'Alasan Terdeteksi',
      'howToDisableFakeGPS': 'Cara Menonaktifkan Fake GPS',
      'fakeGPSStep1':
          'Tutup atau hapus aplikasi Fake GPS atau aplikasi spoofing lokasi di perangkat Anda',
      'fakeGPSStep2':
          'Restart perangkat Anda untuk memastikan semua layanan lokasi direset',
      'fakeGPSStep3': 'Bersihkan cache aplikasi dan izin lokasi',
      'fakeGPSStep4':
          'Coba check-in/check-out lagi setelah menonaktifkan Fake GPS',
      'fakeGPSImportantNote':
          'Fitur check-in dan check-out memerlukan lokasi GPS asli untuk verifikasi. Menggunakan Fake GPS tidak diizinkan dan akan mencegah Anda menggunakan fitur-fitur ini.',
      'close': 'Tutup',
      'retake': 'Ambil Ulang',
      'selectContact': 'Pilih Kontak',
      'selectContactOptional': 'Pilih kontak (opsional)',
      'selectVisitDate': 'Pilih tanggal kunjungan',
      'enterVisitPurpose': 'Masukkan tujuan kunjungan',
      'enterAdditionalNotes': 'Masukkan catatan tambahan',
      'pleaseSelectAccount': 'Harap pilih akun',
      'deal': 'Deal',
      'lead': 'Lead',
      'none': 'Tidak Ada',
      'pleaseSelectDeal': 'Silakan pilih deal',
      'pleaseSelectLead': 'Silakan pilih lead',
      'noLeadsAvailable': 'Tidak ada lead tersedia',
      'purposeRequired': 'Tujuan wajib diisi',
      'updateVisitReport': 'Perbarui Laporan Kunjungan',
      'deleteVisitReport': 'Hapus Laporan Kunjungan',
      'submitVisitReport': 'Kirim Laporan Kunjungan',
      'deleteVisitReportConfirmation':
          'Apakah Anda yakin ingin menghapus laporan kunjungan ini? Tindakan ini tidak dapat dibatalkan.',
      'visitReportUpdatedSuccessfully': 'Laporan kunjungan berhasil diperbarui',
      'visitReportDeletedSuccessfully': 'Laporan kunjungan berhasil dihapus',
      'visitReportSubmittedSuccessfully': 'Laporan kunjungan berhasil dikirim',
      'failedToUpdateVisitReport': 'Gagal memperbarui laporan kunjungan',
      'failedToDeleteVisitReport': 'Gagal menghapus laporan kunjungan',
      'failedToSubmitVisitReport': 'Gagal mengirim laporan kunjungan',
      'outcome': 'Hasil',
      'nextSteps': 'Langkah Selanjutnya',
      'selectOutcome': 'Pilih hasil',
      'enterNextSteps': 'Masukkan langkah selanjutnya (opsional)',
      'visitReportCreatedSuccessfully': 'Laporan kunjungan berhasil dibuat',
      'failedToCreateVisitReport': 'Gagal membuat laporan kunjungan',
      'all': 'Semua',
      // Route Optimization
      'routeOptimization': 'Optimisasi Rute',
      'createRoute': 'Buat Rute',
      'routeDetails': 'Detail Rute',
      'routeName': 'Nama Rute',
      'routeNameOptional': 'Nama Rute (Opsional)',
      'enterRouteName': 'Masukkan nama rute',
      'startLocation': 'Lokasi Awal',
      'notSet': 'Belum diatur',
      'useCurrentLocation': 'Gunakan lokasi saat ini',
      'waypoints': 'Titik Perhentian',
      'add': 'Tambah',
      'noWaypointsAdded': 'Belum ada titik perhentian ditambahkan',
      'optimize': 'Optimalkan',
      'noRoutesFound': 'Tidak ada rute ditemukan',
      'createFirstOptimizedRoute': 'Buat rute teroptimasi pertama Anda',
      'deleteRoute': 'Hapus Rute',
      'deleteRouteConfirmation':
          'Apakah Anda yakin ingin menghapus rute ini? Tindakan ini tidak dapat dibatalkan.',
      'routeDeleted': 'Rute berhasil dihapus',
      'failedToDeleteRoute': 'Gagal menghapus rute',
      'failedToOptimizeRoute': 'Gagal mengoptimalkan rute',
      'routeOptimizedSuccessfully': 'Rute berhasil dioptimalkan',
      'unnamedRoute': 'Rute Tanpa Nama',
      'distance': 'Jarak',
      'duration': 'Durasi',
      'stops': 'Perhentian',
      'routeSteps': 'Langkah Rute',
      'continueRoute': 'Lanjutkan',
      'routeMenu': 'Rute',
      'selectWaypoint': 'Pilih Titik Perhentian',
      'selectFromAccounts': 'Pilih dari Akun',
      'selectFromVisitReports': 'Pilih dari Laporan Kunjungan',
      'searchAccountsOrContacts': 'Cari akun atau kontak...',
      'noAccountsFoundForWaypoint': 'Tidak ada akun ditemukan',
      'noVisitReportsFoundForWaypoint': 'Tidak ada laporan kunjungan ditemukan',
      'destination': 'Tujuan',
      // Leads
      'leads': 'Prospek',
      'pipeline': 'Pipeline',
      'createLead': 'Buat Prospek',
      'editLead': 'Edit Prospek',
      'leadDetails': 'Detail Prospek',
      'searchLeads': 'Cari prospek...',
      'noLeadsFound': 'Tidak ada prospek ditemukan',
      'deleteLead': 'Hapus Prospek',
      'deleteLeadConfirmation':
          'Apakah Anda yakin ingin menghapus prospek ini?',
      'leadDeletedSuccessfully': 'Prospek berhasil dihapus',
      'failedToDeleteLead': 'Gagal menghapus prospek',
      'leadCreatedSuccessfully': 'Prospek berhasil dibuat',
      'leadUpdatedSuccessfully': 'Prospek berhasil diperbarui',
      'firstName': 'Nama Depan',
      'lastName': 'Nama Belakang',
      'failedToCreateLead': 'Gagal membuat prospek',
      'failedToUpdateLead': 'Gagal memperbarui prospek',
      'convertLead': 'Konversi Prospek',
      'convertLeadConfirmation':
          'Apakah Anda yakin ingin mengonversi prospek ini?',
      'leadConvertedSuccessfully': 'Prospek berhasil dikonversi',
      'failedToConvertLead': 'Gagal mengonversi prospek',
      'company': 'Perusahaan',
      'industry': 'Industri',
      'source': 'Sumber',
      'website': 'Website',
      'convert': 'Konversi',
      'opportunityTitle': 'Judul Peluang',
      'pipelineStage': 'Tahapan Pipeline',
      'dealValue': 'Nilai Deal',
      'expectedCloseDate': 'Estimasi Tanggal Closing',
      'probability': 'Probabilitas (%)',
      'isRequired': 'Wajib diisi',
      'clear': 'Bersihkan',
      'apply': 'Terapkan',
      'filterBySource': 'Filter berdasarkan Sumber',
      'filterByIndustry': 'Filter berdasarkan Industri',
      'filterByProvince': 'Filter berdasarkan Provinsi',
      'searchSchedules': 'Cari jadwal...',
      'noSchedulesFound': 'Tidak ada jadwal ditemukan',
      'tapToCreateSchedule': 'Ketuk + untuk membuat jadwal baru',
      'schedules': 'Jadwal',
      'scheduleDetails': 'Detail Jadwal',
      'scheduleInformation': 'Informasi Jadwal',
      'scheduledAt': 'Waktu Jadwal',
      'reminder': 'Pengingat',
      'minutesBefore': 'menit sebelum',
      'linkedTask': 'Tugas Terkait',
      'selectTask': 'Pilih Tugas',
      'scheduledAtRequired': 'Waktu jadwal wajib diisi',
      'scheduledFrom': 'Jadwal Dari',
      'scheduledTo': 'Jadwal Sampai',
      'selectDate': 'Pilih Tanggal',
      'selectScheduledFrom': 'Pilih jadwal dari',
      'selectScheduledTo': 'Pilih jadwal sampai',
      'clearDate': 'Bersihkan Tanggal',
      'noData': 'Tidak ada data',
      'noDescription': 'Tidak ada deskripsi',
      'location': 'Lokasi',
      'noLocation': 'Tidak ada lokasi',
      'relatedVisitReport': 'Laporan Kunjungan Terkait',
      'viewVisitReport': 'Lihat Laporan Kunjungan',
      'activityType': 'Tipe Aktivitas',
      'scheduled': 'Terjadwal',
      'createSchedule': 'Buat Jadwal',
      'editSchedule': 'Edit Jadwal',
      'updateSchedule': 'Perbarui Jadwal',
      'deleteSchedule': 'Hapus Jadwal',
      'deleteScheduleConfirmation':
          'Apakah Anda yakin ingin menghapus jadwal ini? Tindakan ini tidak dapat dibatalkan.',
      'scheduleCreatedSuccessfully': 'Jadwal berhasil dibuat',
      'scheduleUpdatedSuccessfully': 'Jadwal berhasil diperbarui',
      'scheduleDeletedSuccessfully': 'Jadwal berhasil dihapus',
      'failedToCreateSchedule': 'Gagal membuat jadwal',
      'failedToUpdateSchedule': 'Gagal memperbarui jadwal',
      'failedToDeleteSchedule': 'Gagal menghapus jadwal',
      'reminderMinutesBefore': 'Pengingat (menit sebelum)',
      'invalidReminderMinutes': 'Menit pengingat tidak valid (0-10080)',
      'enterTitle': 'Masukkan judul',
      'pending': 'Tertunda',
      'confirmed': 'Dikonfirmasi',
      'minCharacters': 'Minimal {count} karakter',
      'maxCharacters': 'Maksimal {count} karakter',
    },
  };
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) {
    return AppLocalizations.supportedLocales.any(
      (supportedLocale) => supportedLocale.languageCode == locale.languageCode,
    );
  }

  @override
  Future<AppLocalizations> load(Locale locale) async {
    return AppLocalizations(locale);
  }

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}
