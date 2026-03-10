gpu-scheduler-platform/
├── README.md
├── LICENSE
├── Makefile
├── go.mod
├── go.sum
├── .gitignore
├── .golangci.yml
├── buf.yaml
├── buf.gen.yaml
├── codecov.yaml
├── Dockerfile
├── Dockerfile.api
├── Dockerfile.scheduler
├── Dockerfile.controller
├── Dockerfile.webhook
├── Dockerfile.agent
│
├── cmd/
│   ├── api-server/
│   │   └── main.go
│   ├── scheduler/
│   │   └── main.go
│   ├── controller/
│   │   └── main.go
│   ├── webhook/
│   │   └── main.go
│   └── agent/
│       └── main.go
│
├── configs/
│   ├── api-server.yaml
│   ├── scheduler.yaml
│   ├── controller.yaml
│   ├── webhook.yaml
│   ├── agent.yaml
│   ├── policy.sample.yaml
│   ├── logging.sample.yaml
│   └── features.sample.yaml
│
├── api/
│   ├── openapi/
│   │   ├── api-server.yaml
│   │   └── swagger.yaml
│   ├── proto/
│   │   ├── cluster/v1/
│   │   │   └── cluster.proto
│   │   ├── job/v1/
│   │   │   └── job.proto
│   │   ├── scheduler/v1/
│   │   │   └── scheduler.proto
│   │   └── nodeagent/v1/
│   │       └── agent.proto
│   └── crd/
│       ├── gpujobs.yaml
│       ├── gpupolicies.yaml
│       ├── gpuqueues.yaml
│       ├── gpuquotas.yaml
│       └── gpuclustersnapshots.yaml
│
├── deployments/
│   ├── docker/
│   │   ├── api-server/
│   │   ├── scheduler/
│   │   ├── controller/
│   │   ├── webhook/
│   │   └── agent/
│   ├── kustomize/
│   │   ├── base/
│   │   │   ├── namespace.yaml
│   │   │   ├── serviceaccount.yaml
│   │   │   ├── rbac.yaml
│   │   │   ├── configmap.yaml
│   │   │   ├── secret.yaml
│   │   │   ├── api-server-deployment.yaml
│   │   │   ├── scheduler-deployment.yaml
│   │   │   ├── controller-deployment.yaml
│   │   │   ├── webhook-deployment.yaml
│   │   │   ├── agent-daemonset.yaml
│   │   │   ├── service.yaml
│   │   │   ├── webhook-config.yaml
│   │   │   ├── certificate.yaml
│   │   │   ├── crds.yaml
│   │   │   └── kustomization.yaml
│   │   ├── overlays/
│   │   │   ├── dev/
│   │   │   ├── test/
│   │   │   └── prod/
│   │   └── addons/
│   │       ├── prometheus-servicemonitor.yaml
│   │       ├── grafana-dashboard.json
│   │       └── alert-rules.yaml
│   └── helm/
│       └── gpu-scheduler-platform/
│           ├── Chart.yaml
│           ├── values.yaml
│           ├── values-dev.yaml
│           ├── values-prod.yaml
│           ├── templates/
│           │   ├── namespace.yaml
│           │   ├── serviceaccount.yaml
│           │   ├── rbac.yaml
│           │   ├── configmap.yaml
│           │   ├── secret.yaml
│           │   ├── api-server-deployment.yaml
│           │   ├── scheduler-deployment.yaml
│           │   ├── controller-deployment.yaml
│           │   ├── webhook-deployment.yaml
│           │   ├── agent-daemonset.yaml
│           │   ├── services.yaml
│           │   ├── ingress.yaml
│           │   ├── pdb.yaml
│           │   ├── servicemonitor.yaml
│           │   ├── prometheusrule.yaml
│           │   ├── validatingwebhookconfiguration.yaml
│           │   ├── mutatingwebhookconfiguration.yaml
│           │   ├── certificates.yaml
│           │   └── _helpers.tpl
│           └── charts/
│
├── internal/
│   ├── app/
│   │   ├── apiserver/
│   │   │   ├── app.go
│   │   │   ├── routes.go
│   │   │   ├── hooks.go
│   │   │   └── wire.go
│   │   ├── scheduler/
│   │   │   ├── app.go
│   │   │   ├── runner.go
│   │   │   ├── leader.go
│   │   │   └── wire.go
│   │   ├── controller/
│   │   │   ├── app.go
│   │   │   ├── manager.go
│   │   │   ├── controllers.go
│   │   │   └── wire.go
│   │   ├── webhook/
│   │   │   ├── app.go
│   │   │   ├── server.go
│   │   │   └── wire.go
│   │   └── agent/
│   │       ├── app.go
│   │       ├── reporter.go
│   │       └── wire.go
│   │
│   ├── bootstrap/
│   │   ├── config.go
│   │   ├── logger.go
│   │   ├── metrics.go
│   │   ├── tracing.go
│   │   ├── pprof.go
│   │   ├── mysql.go
│   │   ├── redis.go
│   │   ├── k8s.go
│   │   ├── grpc.go
│   │   ├── http.go
│   │   ├── leader_election.go
│   │   └── lifecycle.go
│   │
│   ├── config/
│   │   ├── types.go
│   │   ├── defaults.go
│   │   ├── validate.go
│   │   ├── loader_yaml.go
│   │   ├── loader_env.go
│   │   └── feature_gates.go
│   │
│   ├── domain/
│   │   ├── cluster/
│   │   │   ├── node.go
│   │   │   ├── gpu.go
│   │   │   ├── topology.go
│   │   │   ├── mig.go
│   │   │   └── snapshot.go
│   │   ├── job/
│   │   │   ├── job.go
│   │   │   ├── requirement.go
│   │   │   ├── queue.go
│   │   │   ├── priority.go
│   │   │   └── status.go
│   │   ├── policy/
│   │   │   ├── quota.go
│   │   │   ├── fairness.go
│   │   │   ├── preemption.go
│   │   │   ├── placement.go
│   │   │   └── feature_flag.go
│   │   ├── allocation/
│   │   │   ├── allocation.go
│   │   │   ├── reservation.go
│   │   │   └── binding.go
│   │   └── event/
│   │       ├── event.go
│   │       └── reason.go
│   │
│   ├── scheduler/
│   │   ├── queue/
│   │   │   ├── interface.go
│   │   │   ├── priority_queue.go
│   │   │   ├── fair_queue.go
│   │   │   └── aging.go
│   │   ├── framework/
│   │   │   ├── framework.go
│   │   │   ├── cycle_state.go
│   │   │   ├── plugin.go
│   │   │   ├── registry.go
│   │   │   └── status.go
│   │   ├── plugins/
│   │   │   ├── filter/
│   │   │   │   ├── resource_fit.go
│   │   │   │   ├── model_match.go
│   │   │   │   ├── mig_fit.go
│   │   │   │   ├── topology_fit.go
│   │   │   │   └── tenant_quota.go
│   │   │   ├── score/
│   │   │   │   ├── binpack.go
│   │   │   │   ├── spread.go
│   │   │   │   ├── topology_score.go
│   │   │   │   ├── fragmentation_score.go
│   │   │   │   └── utilization_score.go
│   │   │   ├── reserve/
│   │   │   │   └── reservation.go
│   │   │   ├── permit/
│   │   │   │   └── gang_permit.go
│   │   │   ├── prebind/
│   │   │   │   └── binding_check.go
│   │   │   └── preempt/
│   │   │       ├── victim_search.go
│   │   │       └── reclaim.go
│   │   ├── algorithm/
│   │   │   ├── schedule_one.go
│   │   │   ├── select_node.go
│   │   │   ├── select_gpu.go
│   │   │   ├── reserve.go
│   │   │   └── commit.go
│   │   ├── cache/
│   │   │   ├── snapshot.go
│   │   │   ├── node_cache.go
│   │   │   ├── job_cache.go
│   │   │   └── reservation_cache.go
│   │   ├── service/
│   │   │   ├── scheduler_service.go
│   │   │   ├── placement_service.go
│   │   │   ├── preemption_service.go
│   │   │   └── fairness_service.go
│   │   └── metrics/
│   │       ├── metrics.go
│   │       └── labels.go
│   │
│   ├── controller/
│   │   ├── reconciler/
│   │   │   ├── gpujob_controller.go
│   │   │   ├── gpupolicy_controller.go
│   │   │   ├── gpuqueue_controller.go
│   │   │   ├── gpuquota_controller.go
│   │   │   ├── pod_controller.go
│   │   │   ├── node_controller.go
│   │   │   └── lease_controller.go
│   │   ├── watcher/
│   │   │   ├── node_watcher.go
│   │   │   ├── pod_watcher.go
│   │   │   ├── job_watcher.go
│   │   │   └── event_watcher.go
│   │   ├── syncer/
│   │   │   ├── cluster_syncer.go
│   │   │   ├── allocation_syncer.go
│   │   │   ├── job_status_syncer.go
│   │   │   └── metrics_syncer.go
│   │   └── service/
│   │       ├── controller_service.go
│   │       └── snapshot_service.go
│   │
│   ├── webhook/
│   │   ├── mutating/
│   │   │   ├── gpujob_defaults.go
│   │   │   └── pod_defaults.go
│   │   ├── validating/
│   │   │   ├── gpujob_validate.go
│   │   │   ├── gpupolicy_validate.go
│   │   │   ├── gpuquota_validate.go
│   │   │   └── pod_validate.go
│   │   └── server/
│   │       ├── router.go
│   │       └── handlers.go
│   │
│   ├── apiserver/
│   │   ├── handler/
│   │   │   ├── health_handler.go
│   │   │   ├── ready_handler.go
│   │   │   ├── job_handler.go
│   │   │   ├── queue_handler.go
│   │   │   ├── policy_handler.go
│   │   │   ├── quota_handler.go
│   │   │   ├── cluster_handler.go
│   │   │   ├── tenant_handler.go
│   │   │   └── metrics_handler.go
│   │   ├── service/
│   │   │   ├── job_service.go
│   │   │   ├── policy_service.go
│   │   │   ├── quota_service.go
│   │   │   ├── cluster_service.go
│   │   │   └── tenant_service.go
│   │   ├── dto/
│   │   │   ├── common.go
│   │   │   ├── job.go
│   │   │   ├── queue.go
│   │   │   ├── policy.go
│   │   │   ├── quota.go
│   │   │   ├── cluster.go
│   │   │   └── tenant.go
│   │   └── router/
│   │       ├── routes.go
│   │       └── middleware.go
│   │
│   ├── agent/
│   │   ├── collector/
│   │   │   ├── nvidia_smi.go
│   │   │   ├── dcgm.go
│   │   │   ├── mig.go
│   │   │   ├── topology.go
│   │   │   └── pod_gpu_usage.go
│   │   ├── reporter/
│   │   │   ├── grpc_reporter.go
│   │   │   ├── http_reporter.go
│   │   │   └── heartbeat.go
│   │   ├── discovery/
│   │   │   ├── device_discovery.go
│   │   │   └── node_meta.go
│   │   └── service/
│   │       └── agent_service.go
│   │
│   ├── repo/
│   │   ├── models/
│   │   │   ├── gpu_job.go
│   │   │   ├── gpu_queue.go
│   │   │   ├── gpu_quota.go
│   │   │   ├── gpu_policy.go
│   │   │   ├── node_snapshot.go
│   │   │   ├── allocation.go
│   │   │   ├── tenant.go
│   │   │   ├── audit_log.go
│   │   │   └── outbox.go
│   │   ├── mysql/
│   │   │   ├── gpu_job_repo.go
│   │   │   ├── gpu_queue_repo.go
│   │   │   ├── gpu_quota_repo.go
│   │   │   ├── gpu_policy_repo.go
│   │   │   ├── node_snapshot_repo.go
│   │   │   ├── allocation_repo.go
│   │   │   ├── tenant_repo.go
│   │   │   ├── audit_log_repo.go
│   │   │   └── tx.go
│   │   ├── redis/
│   │   │   ├── keys.go
│   │   │   ├── queue_cache.go
│   │   │   ├── snapshot_cache.go
│   │   │   ├── reservation_cache.go
│   │   │   ├── limiter.go
│   │   │   └── lock.go
│   │   └── interface.go
│   │
│   ├── k8s/
│   │   ├── client/
│   │   │   ├── clientset.go
│   │   │   ├── informer.go
│   │   │   └── dynamic.go
│   │   ├── crd/
│   │   │   ├── types.go
│   │   │   ├── zz_generated.deepcopy.go
│   │   │   └── register.go
│   │   ├── scheduler/
│   │   │   ├── extender.go
│   │   │   ├── framework_plugin.go
│   │   │   └── binder.go
│   │   ├── pod/
│   │   │   ├── patcher.go
│   │   │   ├── annotations.go
│   │   │   └── finalizer.go
│   │   └── node/
│   │       ├── labels.go
│   │       └── taints.go
│   │
│   ├── middleware/
│   │   ├── request_id.go
│   │   ├── recovery.go
│   │   ├── access_log.go
│   │   ├── authn.go
│   │   ├── authz.go
│   │   ├── ratelimit.go
│   │   └── audit.go
│   │
│   ├── observability/
│   │   ├── metrics/
│   │   │   ├── scheduler.go
│   │   │   ├── controller.go
│   │   │   ├── api.go
│   │   │   ├── agent.go
│   │   │   └── registry.go
│   │   ├── tracing/
│   │   │   └── otel.go
│   │   ├── logging/
│   │   │   ├── logger.go
│   │   │   ├── fields.go
│   │   │   └── audit.go
│   │   └── profiling/
│   │       └── pprof.go
│   │
│   ├── auth/
│   │   ├── jwt.go
│   │   ├── rbac.go
│   │   ├── subject.go
│   │   └── permission.go
│   │
│   ├── eventbus/
│   │   ├── publisher.go
│   │   ├── consumer.go
│   │   ├── topics.go
│   │   └── outbox.go
│   │
│   ├── audit/
│   │   ├── audit.go
│   │   ├── recorder.go
│   │   └── formatter.go
│   │
│   └── util/
│       ├── pointer.go
│       ├── time.go
│       ├── uuid.go
│       ├── errors.go
│       ├── retry.go
│       ├── backoff.go
│       ├── json.go
│       ├── yaml.go
│       └── net.go
│
├── pkg/
│   ├── client/
│   │   ├── apiclient/
│   │   └── schedulerclient/
│   ├── types/
│   │   ├── labels.go
│   │   ├── annotations.go
│   │   └── constants.go
│   └── version/
│       └── version.go
│
├── database/
│   ├── migrations/
│   │   ├── 000001_init.up.sql
│   │   ├── 000001_init.down.sql
│   │   ├── 000002_add_gpu_queue.up.sql
│   │   ├── 000002_add_gpu_queue.down.sql
│   │   ├── 000003_add_reservation.up.sql
│   │   ├── 000003_add_reservation.down.sql
│   │   └── ...
│   ├── seeds/
│   │   ├── dev_seed.sql
│   │   └── prod_seed.sql
│   └── diagrams/
│       └── er.png
│
├── scripts/
│   ├── dev/
│   │   ├── run-local.sh
│   │   ├── run-api.sh
│   │   ├── run-scheduler.sh
│   │   ├── run-controller.sh
│   │   ├── run-webhook.sh
│   │   └── run-agent.sh
│   ├── build/
│   │   ├── build-images.sh
│   │   ├── build-binaries.sh
│   │   └── lint.sh
│   ├── ci/
│   │   ├── test.sh
│   │   ├── integration.sh
│   │   └── e2e.sh
│   ├── db/
│   │   ├── migrate-up.sh
│   │   ├── migrate-down.sh
│   │   └── seed.sh
│   ├── k8s/
│   │   ├── install-crd.sh
│   │   ├── deploy-dev.sh
│   │   ├── deploy-prod.sh
│   │   └── uninstall.sh
│   └── tools/
│       ├── gen-crd.sh
│       ├── gen-openapi.sh
│       └── gen-proto.sh
│
├── test/
│   ├── e2e/
│   │   ├── suite_test.go
│   │   ├── scheduler_e2e_test.go
│   │   ├── quota_e2e_test.go
│   │   ├── preemption_e2e_test.go
│   │   └── fixtures/
│   ├── integration/
│   │   ├── mysql/
│   │   ├── redis/
│   │   ├── k8s/
│   │   └── scheduler/
│   ├── performance/
│   │   ├── benchmark_schedule_test.go
│   │   └── benchmark_cache_test.go
│   └── fuzz/
│       └── queue_fuzz_test.go
│
├── docs/
│   ├── architecture/
│   │   ├── overview.md
│   │   ├── scheduler-framework.md
│   │   ├── queue-design.md
│   │   ├── preemption.md
│   │   ├── mig-support.md
│   │   ├── topology-aware.md
│   │   └── deployment.md
│   ├── api/
│   │   └── examples.md
│   ├── runbooks/
│   │   ├── troubleshooting.md
│   │   ├── rollback.md
│   │   ├── cert-rotation.md
│   │   └── disaster-recovery.md
│   └── adr/
│       ├── 0001-use-crd-for-job-model.md
│       ├── 0002-use-mysql-plus-redis.md
│       └── 0003-plugin-based-scheduler-framework.md
│
└── hack/
├── boilerplate.go.txt
├── tools.go
└── kind/
├── kind-cluster.yaml
└── gpu-node-mock.yaml