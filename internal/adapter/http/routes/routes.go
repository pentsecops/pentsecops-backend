package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/adapter/http/handlers"
	middleware_pkg "github.com/pentsecops/backend/pkg/middleware"
)

// SetupRoutes sets up all application routes
func SetupRoutes(
	app *fiber.App,
	authHandler *handlers.AuthHandler,
	authMiddleware *middleware_pkg.AuthMiddleware,
	adminOverviewHandler *handlers.AdminOverviewHandler,
	usersHandler *handlers.UsersHandler,
	projectsHandler *handlers.ProjectsHandler,
	tasksHandler *handlers.TasksHandler,
	vulnerabilitiesHandler *handlers.VulnerabilitiesHandler,
	domainsHandler *handlers.DomainsHandler,
	notificationsHandler *handlers.NotificationsHandler,
	auditHandler *handlers.AuditHandler,
	adminNotificationsHandler *handlers.AdminNotificationsHandler,
	pentesterOverviewHandler *handlers.PentesterOverviewHandler,
	pentesterProjectsHandler *handlers.PentesterProjectsHandler,
	pentesterTasksHandler *handlers.PentesterTasksHandler,
	pentesterSubmitReportHandler *handlers.PentesterSubmitReportHandler,
	pentesterAlertsHandler *handlers.PentesterAlertsHandler,
	stakeholderOverviewHandler *handlers.StakeholderOverviewHandler,
	stakeholderVulnerabilitiesHandler *handlers.StakeholderVulnerabilitiesHandler,
	stakeholderReportsHandler *handlers.StakeholderReportsHandler,
	llmHandler *handlers.LLMHandler,
) {
	// Add activity logger middleware globally (will be set up in main.go with auditRepo)

	// LLM Query endpoint (protected, all roles)
	api := app.Group("/api")
	apiV1 := api.Group("/v1", authMiddleware.RequireAuth())
	apiV1.Post("/query", llmHandler.Query)

	// Health check endpoint (no auth required)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "PentSecOps API is running",
		})
	})

	// Test route to verify admin routes are working
	api.Get("/test-admin", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Admin routes are accessible",
		})
	})

	// Authentication routes (no auth required, but with rate limiting)
	auth := api.Group("/auth", middleware_pkg.AuthRateLimit())
	{
		auth.Post("/login", authHandler.Login)                                                  // Login
		auth.Post("/refresh", authHandler.RefreshToken)                                         // Refresh token
		auth.Post("/logout", authHandler.Logout)                                                // Logout
		auth.Post("/change-password", authMiddleware.RequireAuth(), authHandler.ChangePassword) // Change password (requires auth)
	}

	// Admin routes (requires admin role)
	admin := api.Group("/admin", authMiddleware.RequireRole("admin"))
	{
		// Overview routes
		overview := admin.Group("/overview")
		{
			overview.Get("/stats", adminOverviewHandler.GetOverviewStats)
			overview.Get("/vulnerabilities-by-severity", adminOverviewHandler.GetVulnerabilitiesBySeverity)
			overview.Get("/top-domains", adminOverviewHandler.GetTop5Domains)
			overview.Get("/project-status", adminOverviewHandler.GetProjectStatusDistribution)
			overview.Get("/recent-activity", adminOverviewHandler.GetRecentActivity)
		}

		// Users routes
		users := admin.Group("/users")
		{
			users.Get("/", usersHandler.ListUsers)              // UC17: List users with pagination
			users.Get("/refresh", usersHandler.RefreshUsers)    // UC12: Refresh users list
			users.Post("/", usersHandler.CreateUser)            // UC14: Create new user
			users.Put("/:id", usersHandler.UpdateUser)          // Update user
			users.Delete("/:id", usersHandler.DeleteUser)       // UC19: Delete user
			users.Get("/stats", usersHandler.GetUserStats)      // UC24: Get user statistics
			users.Get("/export", usersHandler.ExportUsersToCSV) // UC25: Export users to CSV
		}

		// Projects routes
		projects := admin.Group("/projects")
		{
			// Projects sub-tab
			projects.Get("/", projectsHandler.ListProjects)            // UC31, UC32: List projects with pagination
			projects.Post("/", projectsHandler.CreateProject)          // UC28: Create new project
			projects.Put("/:id", projectsHandler.UpdateProject)        // Update project
			projects.Delete("/:id", projectsHandler.DeleteProject)     // Delete project
			projects.Get("/stats", projectsHandler.GetProjectStats)    // UC27: Get project statistics
			projects.Get("/pentesters", projectsHandler.GetPentesters) // Get pentesters for dropdown

			// Task Board sub-tab
			projects.Get("/tasks", tasksHandler.ListAllTasks)                           // UC36: Display task board
			projects.Get("/tasks/project/:project_id", tasksHandler.ListTasksByProject) // UC36: List tasks by project
			projects.Post("/tasks", tasksHandler.CreateTask)                            // UC37: Create new task
			projects.Patch("/tasks/:id/status", tasksHandler.UpdateTaskStatus)          // UC38: Update task status
			projects.Delete("/tasks/:id", tasksHandler.DeleteTask)                      // Delete task
		}

		// Vulnerabilities routes
		vulnerabilities := admin.Group("/vulnerabilities")
		{
			vulnerabilities.Get("/", vulnerabilitiesHandler.ListVulnerabilities)              // UC40-43, UC48-49: List vulnerabilities with search and filters
			vulnerabilities.Post("/", vulnerabilitiesHandler.CreateVulnerability)             // UC44-45: Create new vulnerability
			vulnerabilities.Get("/stats", vulnerabilitiesHandler.GetVulnerabilityStats)       // UC39: Get vulnerability statistics
			vulnerabilities.Get("/sla", vulnerabilitiesHandler.GetSLACompliance)              // UC54: Get SLA compliance
			vulnerabilities.Get("/export", vulnerabilitiesHandler.ExportVulnerabilitiesToCSV) // UC53: Export vulnerabilities to CSV
			vulnerabilities.Get("/:id", vulnerabilitiesHandler.GetVulnerabilityByID)          // Get vulnerability by ID
			vulnerabilities.Put("/:id", vulnerabilitiesHandler.UpdateVulnerability)           // UC52: Update vulnerability
			vulnerabilities.Delete("/:id", vulnerabilitiesHandler.DeleteVulnerability)        // Delete vulnerability
		}

		// Domains routes
		domains := admin.Group("/domains")
		{
			domains.Get("/", domainsHandler.ListDomains)                           // UC59-60: List domains with pagination
			domains.Post("/", domainsHandler.CreateDomain)                         // Create new domain
			domains.Get("/stats", domainsHandler.GetDomainsStats)                  // UC55-58: Get domains statistics
			domains.Get("/security-metrics", domainsHandler.GetSecurityMetrics)    // UC65: Get security metrics for radar chart
			domains.Post("/security-metrics", domainsHandler.CreateSecurityMetric) // Create security metric
			domains.Get("/sla-breach", domainsHandler.GetSLABreachAnalysis)        // UC66: Get SLA breach analysis
			domains.Get("/:id", domainsHandler.GetDomainByID)                      // Get domain by ID
			domains.Put("/:id", domainsHandler.UpdateDomain)                       // Update domain
			domains.Delete("/:id", domainsHandler.DeleteDomain)                    // Delete domain
		}

		// Notifications routes
		notifications := admin.Group("/notifications")
		{
			notifications.Get("/total", notificationsHandler.GetTotalNotificationsSent) // UC81: Get total notifications sent
			notifications.Get("/", notificationsHandler.ListNotifications)              // UC82, UC83: List notifications with pagination
			notifications.Post("/", notificationsHandler.CreateNotification)            // UC79: Create new notification
			notifications.Get("/alerts", notificationsHandler.ListImportantAlerts)      // UC84, UC85: List important alerts from pentesters
		}

		// Direct Tasks routes (for admin task management)
		tasks := admin.Group("/tasks")
		{
			tasks.Get("/", tasksHandler.ListAllTasks)                          // UC36: Display task board
			tasks.Get("/project/:project_id", tasksHandler.ListTasksByProject) // UC36: List tasks by project
			tasks.Post("/", tasksHandler.CreateTask)                           // UC37: Create new task
			tasks.Patch("/:id/status", tasksHandler.UpdateTaskStatus)          // UC38: Update task status
			tasks.Delete("/:id", tasksHandler.DeleteTask)                      // Delete task
		}

		// Audit Logs routes (for admin activity monitoring)
		audit := admin.Group("/audit")
		{
			audit.Get("/logs", auditHandler.GetActivityLogs)      // Get activity logs with filters
			audit.Get("/stats", auditHandler.GetActivityStats)    // Get activity statistics
			audit.Get("/export", auditHandler.ExportActivityLogs) // Export activity logs to CSV
		}

		// Admin Notifications routes (for sending notifications)
		adminNotifications := admin.Group("/notifications")
		{
			adminNotifications.Post("/send", adminNotificationsHandler.SendNotification)     // Send notification
			adminNotifications.Get("/", adminNotificationsHandler.GetNotifications)          // Get sent notifications
			adminNotifications.Get("/stats", adminNotificationsHandler.GetNotificationStats) // Get notification statistics
			adminNotifications.Get("/users", adminNotificationsHandler.GetAvailableUsers)    // Get available users for notifications
		}
	}

	// Pentester routes (requires pentester role)
	pentester := api.Group("/pentester", authMiddleware.RequireRole("pentester"))
	{
		// Overview routes
		overview := pentester.Group("/overview")
		{
			overview.Get("/stats", pentesterOverviewHandler.GetOverviewStats)
			overview.Get("/active-projects", pentesterOverviewHandler.GetActiveProjects)
			overview.Get("/recent-vulnerabilities", pentesterOverviewHandler.GetRecentVulnerabilities)
			overview.Get("/upcoming-deadlines", pentesterOverviewHandler.GetUpcomingDeadlines)
		}

		// Projects routes
		projects := pentester.Group("/projects")
		{
			projects.Get("/", pentesterProjectsHandler.GetAssignedProjects)                    // UC9: Fetch and Display Assigned Projects List
			projects.Get("/:id", pentesterProjectsHandler.GetProjectDetails)                   // Get single project details
			projects.Get("/:id/assets", pentesterProjectsHandler.GetProjectAssets)             // UC12: Display Project Assets List
			projects.Get("/:id/requirements", pentesterProjectsHandler.GetProjectRequirements) // UC13: Display Project Requirements List
		}

		// Tasks routes
		tasks := pentester.Group("/tasks")
		{
			tasks.Get("/board", pentesterTasksHandler.GetTaskBoard)              // UC16: Fetch and Display Task Board
			tasks.Get("/projects", pentesterTasksHandler.GetProjectsForDropdown) // Get projects for dropdown
			tasks.Get("/:id", pentesterTasksHandler.GetTaskDetails)              // Get single task details
			tasks.Post("/", pentesterTasksHandler.CreateTask)                    // UC17: Create New Task
			tasks.Patch("/:id/status", pentesterTasksHandler.UpdateTaskStatus)   // UC18: Update Task Status
			tasks.Put("/:id", pentesterTasksHandler.UpdateTask)                  // UC19: Edit Task Details
			tasks.Delete("/:id", pentesterTasksHandler.DeleteTask)               // UC20: Delete Task
		}

		// Submit Report routes
		submitReport := pentester.Group("/submit-report")
		{
			submitReport.Get("/projects", pentesterSubmitReportHandler.GetProjectsForDropdown)    // UC22: Get projects for report dropdown
			submitReport.Post("/", pentesterSubmitReportHandler.SubmitReport)                     // UC23: Submit vulnerability report
			submitReport.Get("/history", pentesterSubmitReportHandler.GetSubmittedReportsHistory) // UC27: Get submitted reports history
			submitReport.Get("/:id", pentesterSubmitReportHandler.GetReportDetails)               // UC31: Get report details
			submitReport.Post("/:id/resubmit", pentesterSubmitReportHandler.ResubmitReport)       // UC30: Resubmit rejected report
		}

		// Alerts routes
		alerts := pentester.Group("/alerts")
		{
			alerts.Get("/statistics", pentesterAlertsHandler.GetAlertStats)          // UC33: Get alert statistics
			alerts.Get("/", pentesterAlertsHandler.GetAlerts)                        // UC34: Get all alerts (with filters)
			alerts.Get("/types", pentesterAlertsHandler.GetAlertTypes)               // UC35: Get alert types
			alerts.Get("/guidelines", pentesterAlertsHandler.GetAlertGuidelines)     // UC39: Get alert guidelines
			alerts.Patch("/:id/read", pentesterAlertsHandler.MarkAlertAsRead)        // UC36: Mark alert as read/unread
			alerts.Delete("/:id", pentesterAlertsHandler.DismissAlert)               // UC37: Dismiss alert
			alerts.Post("/send-to-admin", pentesterAlertsHandler.CreateAlertToAdmin) // UC38: Create alert to admin
		}
	}

	// Stakeholder routes (requires stakeholder role)
	stakeholder := api.Group("/stakeholder")
	stakeholder.Use(authMiddleware.RequireRole("stakeholder"))

	// Overview routes
	stakeholderOverview := stakeholder.Group("/overview")
	stakeholderOverview.Get("/security-metrics", stakeholderOverviewHandler.GetSecurityMetrics)       // UC1-UC6: Security metrics cards
	stakeholderOverview.Get("/vulnerability-trend", stakeholderOverviewHandler.GetVulnerabilityTrend) // UC7: Vulnerability trend chart
	stakeholderOverview.Get("/asset-status", stakeholderOverviewHandler.GetAssetStatus)               // UC8: Asset status chart
	stakeholderOverview.Get("/recent-events", stakeholderOverviewHandler.GetRecentSecurityEvents)     // UC9: Recent security events
	stakeholderOverview.Get("/remediation-updates", stakeholderOverviewHandler.GetRemediationUpdates) // UC10: Remediation updates

	// Vulnerabilities routes
	stakeholderVulnerabilities := stakeholder.Group("/vulnerabilities")
	stakeholderVulnerabilities.Get("/stats", stakeholderVulnerabilitiesHandler.GetVulnerabilitiesStats)         // UC11-UC14: Vulnerabilities statistics
	stakeholderVulnerabilities.Get("/list", stakeholderVulnerabilitiesHandler.ListVulnerabilities)              // UC15-UC23: List vulnerabilities with filters
	stakeholderVulnerabilities.Get("/export-csv", stakeholderVulnerabilitiesHandler.ExportVulnerabilitiesToCSV) // UC24: Export to CSV
	stakeholderVulnerabilities.Get("/sla-compliance", stakeholderVulnerabilitiesHandler.GetSLACompliance)       // UC25-UC27: SLA compliance

	// Reports routes
	stakeholderReports := stakeholder.Group("/reports")
	stakeholderReports.Get("/stats", stakeholderReportsHandler.GetReportsStats)                  // UC28-UC30: Reports statistics
	stakeholderReports.Get("/list", stakeholderReportsHandler.ListReports)                       // UC31-UC36: List reports with filters
	stakeholderReports.Get("/view", stakeholderReportsHandler.ViewReport)                        // UC37-UC38: View report details
	stakeholderReports.Get("/evidence", stakeholderReportsHandler.GetReportEvidenceFiles)        // UC39-UC40: Get evidence files
	stakeholderReports.Get("/download-evidence", stakeholderReportsHandler.DownloadEvidenceFile) // UC41: Download evidence file
	stakeholderReports.Get("/download", stakeholderReportsHandler.DownloadReport)                // UC42: Download report

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "Route not found",
			},
		})
	})
}
