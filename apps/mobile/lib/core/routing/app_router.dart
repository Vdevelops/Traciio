import 'package:flutter/material.dart';

import '../../features/accounts/presentation/account_detail_screen.dart';
import '../../features/accounts/presentation/accounts_screen.dart';
import '../../features/auth/presentation/login_screen.dart';
import '../../features/contacts/presentation/contact_detail_screen.dart';
import '../../features/contacts/presentation/contact_list_screen.dart';
import '../../features/dashboard/presentation/dashboard_screen.dart';
import '../../features/leads/presentation/lead_detail_screen.dart';
import '../../features/leads/presentation/lead_form_screen.dart';
import '../../features/leads/presentation/lead_list_screen.dart';
import '../../features/leads/presentation/lead_convert_screen.dart';
import '../../features/leads/data/models/lead_model.dart';
import '../../features/pipeline/presentation/screens/pipeline_screen.dart';
import '../../features/pipeline/presentation/screens/deal_detail_screen.dart';
import '../../features/pipeline/presentation/screens/deal_form_screen.dart';
import '../../features/profile/presentation/profile_screen.dart';
import '../../features/tasks/presentation/task_detail_screen.dart';
import '../../features/tasks/presentation/task_list_screen.dart';
import '../../features/visit_reports/presentation/reports_screen.dart';
import '../../features/visit_reports/presentation/visit_report_detail_screen.dart';
import '../../features/visit_reports/presentation/simplified_visit_report_form_screen.dart';
import '../../features/notifications/presentation/notification_list_screen.dart';
import '../../features/route_optimization/presentation/route_list_screen.dart';
import '../../features/route_optimization/presentation/route_detail_screen.dart';
import '../../features/schedules/presentation/screens/schedule_list_screen.dart';
import '../../features/schedules/presentation/screens/schedule_detail_screen.dart';
import '../../features/schedules/presentation/screens/schedule_form_screen.dart';
import '../widgets/auth_gate.dart';

class AppRoutes {
  const AppRoutes._();

  static const login = '/auth/login';
  static const dashboard = '/dashboard';
  static const accounts = '/accounts';
  static const contacts = '/contacts';
  static const leads = '/leads';
  static const leadsForm = '/leads/form';
  static const leadsDetail = '/leads/detail';
  static const leadsConvert = '/leads/convert';
  static const visitReports = '/visit-reports';
  static const visitReportsCreate = '/visit-reports/create';
  static const tasks = '/tasks';
  static const tasksCreate = '/tasks/create';
  static const tasksEdit = '/tasks/edit';
  static const notifications = '/notifications';
  static const profile = '/profile';
  static const recentActivities = '/recent-activities';
  static const routeOptimization = '/route-optimization';
  static const pipeline = '/pipeline';
  static const schedules = '/schedules';
  static const scheduleDetail = '/schedules/detail';
  static const scheduleForm = '/schedules/form';
}

class AppRouter {
  const AppRouter._();

  static String get initialRoute => AppRoutes.login;

  static Map<String, WidgetBuilder> get routes => {
    AppRoutes.login: (_) => const LoginScreen(),
    AppRoutes.dashboard: (_) => const AuthGate(
      requiredRoute: AppRoutes.dashboard,
      child: DashboardScreen(),
    ),
    AppRoutes.accounts: (_) => const AuthGate(
      requiredRoute: AppRoutes.accounts,
      child: AccountsScreen(),
    ),
    AppRoutes.leads: (_) =>
        const AuthGate(requiredRoute: AppRoutes.leads, child: LeadListScreen()),
    AppRoutes.leadsForm: (context) {
      final args = ModalRoute.of(context)?.settings.arguments;
      final leadId = args is String ? args : null;
      return AuthGate(
        requiredRoute: AppRoutes.leads,
        child: LeadFormScreen(leadId: leadId),
      );
    },
    AppRoutes.leadsDetail: (context) {
      final args = ModalRoute.of(context)?.settings.arguments;
      final leadId = args is String ? args : '';
      return AuthGate(
        requiredRoute: AppRoutes.leads,
        child: LeadDetailScreen(leadId: leadId),
      );
    },
    AppRoutes.leadsConvert: (context) {
      final lead = ModalRoute.of(context)?.settings.arguments as Lead;
      return AuthGate(
        requiredRoute: AppRoutes.leads,
        child: LeadConvertScreen(lead: lead),
      );
    },
    AppRoutes.profile: (_) => const AuthGate(
      requiredRoute: AppRoutes.profile,
      child: ProfileScreen(),
    ),
    AppRoutes.contacts: (context) {
      final args = ModalRoute.of(context)?.settings.arguments;
      final accountId = args is Map ? args['accountId'] as String? : null;
      return AuthGate(
        requiredRoute: AppRoutes.contacts,
        requiredPermission: 'accounts.view', // Contacts use accounts permission
        child: ContactListScreen(accountId: accountId),
      );
    },
    AppRoutes.visitReports: (_) => const AuthGate(
      requiredRoute: AppRoutes.visitReports,
      child: ReportsScreen(),
    ),
    AppRoutes.visitReportsCreate: (_) => const AuthGate(
      requiredRoute: AppRoutes.visitReportsCreate,
      child: SimplifiedVisitReportFormScreen(),
    ),
    AppRoutes.tasks: (_) =>
        const AuthGate(requiredRoute: AppRoutes.tasks, child: TaskListScreen()),
    // TaskFormScreen routes removed - sales users don't need create/edit tasks
    AppRoutes.notifications: (_) => const AuthGate(
      requiredRoute: AppRoutes.notifications,
      child: NotificationListScreen(),
    ),
    // AppRoutes.recentActivities: (_) => const AuthGate(
    //   requiredRoute: AppRoutes.recentActivities,
    //   child: RecentActivitiesScreen(),
    // ),
    AppRoutes.routeOptimization: (_) => const AuthGate(
      requiredRoute: AppRoutes.routeOptimization,
      child: RouteListScreen(),
    ),
    AppRoutes.pipeline: (_) => const AuthGate(
      requiredRoute: AppRoutes.pipeline,
      child: PipelineScreen(),
    ),
    AppRoutes.schedules: (_) => const AuthGate(
      requiredRoute: AppRoutes.schedules,
      child: ScheduleListScreen(),
    ),
    AppRoutes.scheduleForm: (context) {
      final args = ModalRoute.of(context)?.settings.arguments;
      final scheduleId = args is String ? args : null;
      return AuthGate(
        requiredRoute: AppRoutes.schedules,
        child: ScheduleFormScreen(scheduleId: scheduleId),
      );
    },
    AppRoutes.scheduleDetail: (context) {
      final args = ModalRoute.of(context)?.settings.arguments;
      final scheduleId = args is String ? args : '';
      return AuthGate(
        requiredRoute: AppRoutes.schedules,
        child: ScheduleDetailScreen(scheduleId: scheduleId),
      );
    },
  };

  static Route<dynamic>? onGenerateRoute(RouteSettings settings) {
    final uri = Uri.parse(settings.name ?? '');
    final pathSegments = uri.pathSegments;

    // Account Detail: /accounts/:id
    if (pathSegments.length == 2 &&
        pathSegments[0] == 'accounts' &&
        pathSegments[1].isNotEmpty) {
      return MaterialPageRoute(
        settings: settings,
        builder: (_) => AuthGate(
          requiredRoute: AppRoutes.accounts,
          child: AccountDetailScreen(accountId: pathSegments[1]),
        ),
      );
    }

    // Contact Detail: /contacts/:id
    if (pathSegments.length == 2 &&
        pathSegments[0] == 'contacts' &&
        pathSegments[1].isNotEmpty) {
      return MaterialPageRoute(
        settings: settings,
        builder: (_) => AuthGate(
          requiredRoute: AppRoutes.contacts,
          requiredPermission:
              'accounts.view', // Contacts use accounts permission
          child: ContactDetailScreen(contactId: pathSegments[1]),
        ),
      );
    }

    // Visit Report Detail: /visit-reports/:id
    if (pathSegments.length == 2 &&
        pathSegments[0] == 'visit-reports' &&
        pathSegments[1].isNotEmpty &&
        pathSegments[1] != 'create') {
      return MaterialPageRoute(
        settings: settings,
        builder: (_) => AuthGate(
          requiredRoute: AppRoutes.visitReports,
          child: VisitReportDetailScreen(visitReportId: pathSegments[1]),
        ),
      );
    }

    // Visit Report Form: /visit-reports/create
    if (pathSegments.length == 2 &&
        pathSegments[0] == 'visit-reports' &&
        pathSegments[1] == 'create') {
      return MaterialPageRoute(
        settings: settings,
        builder: (_) => AuthGate(
          requiredRoute: AppRoutes.visitReportsCreate,
          child: const SimplifiedVisitReportFormScreen(),
        ),
      );
    }

    // Task Detail: /tasks/:id
    if (pathSegments.length == 2 &&
        pathSegments[0] == 'tasks' &&
        pathSegments[1].isNotEmpty &&
        pathSegments[1] != 'create' &&
        pathSegments[1] != 'edit') {
      return MaterialPageRoute(
        settings: settings,
        builder: (_) => AuthGate(
          requiredRoute: AppRoutes.tasks,
          child: TaskDetailScreen(taskId: pathSegments[1]),
        ),
      );
    }

    // Task Form routes removed - sales users don't need create/edit tasks

    // Route Optimization Detail: /route-optimization/:id
    if (pathSegments.length == 2 &&
        pathSegments[0] == 'route-optimization' &&
        pathSegments[1].isNotEmpty) {
      return MaterialPageRoute(
        settings: settings,
        builder: (_) => AuthGate(
          requiredRoute: AppRoutes.routeOptimization,
          child: RouteDetailScreen(routeId: pathSegments[1]),
        ),
      );
    }

    // Pipeline Detail/Edit/Create
    if (pathSegments.length >= 2 && pathSegments[0] == 'pipeline') {
      final subPath = pathSegments[1];

      if (subPath == 'create') {
        return MaterialPageRoute(
          settings: settings,
          builder: (_) => const AuthGate(
            requiredRoute: AppRoutes.pipeline,
            child: DealFormScreen(),
          ),
        );
      }

      // /pipeline/edit/:id
      if (subPath == 'edit' && pathSegments.length == 3) {
        return MaterialPageRoute(
          settings: settings,
          builder: (_) => AuthGate(
            requiredRoute: AppRoutes.pipeline,
            child: DealFormScreen(dealId: pathSegments[2]),
          ),
        );
      }

      // /pipeline/:id (Detail)
      if (pathSegments.length == 2) {
        return MaterialPageRoute(
          settings: settings,
          builder: (_) => AuthGate(
            requiredRoute: AppRoutes.pipeline,
            child: DealDetailScreen(dealId: subPath),
          ),
        );
      }
    }

    return null;
  }
}
