SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- tenants
INSERT INTO tenants (id, name, enabled, description)
VALUES
    ('t1', 'tenant-1', 1, 'default tenant for local development')
    ON DUPLICATE KEY UPDATE
                         name = VALUES(name),
                         enabled = VALUES(enabled),
                         description = VALUES(description);

-- node snapshot version v1
INSERT INTO node_snapshots (
    version,
    node_name,
    node_state,
    schedulable,
    labels_json,
    annotations_json,
    topology_json,
    created_at
)
VALUES
    (
        'v1',
        'node-a',
        'ready',
        1,
        JSON_OBJECT('gpu.vendor', 'nvidia', 'zone', 'az-a'),
        JSON_OBJECT(),
        JSON_OBJECT(
                'labels', JSON_OBJECT('gpu.vendor', 'nvidia', 'zone', 'az-a'),
                'annotations', JSON_OBJECT(),
                'topology', JSON_OBJECT(
                        'NodeName', 'node-a',
                        'Links', JSON_ARRAY(
                                JSON_OBJECT('FromGPU', 'gpu-a0', 'ToGPU', 'gpu-a1', 'Type', 'nvlink', 'Weight', 10)
                                 )
                            ),
                'gpus', JSON_ARRAY(
                        JSON_OBJECT(
                                'ID', 'gpu-a0',
                                'UUID', 'GPU-A0-UUID',
                                'NodeName', 'node-a',
                                'Index', 0,
                                'Model', 'A100',
                                'Vendor', 'nvidia',
                                'Type', 'full',
                                'MemoryMiB', 81920,
                                'FreeMemoryMiB', 81920,
                                'Healthy', true,
                                'Health', 'healthy',
                                'MIGEnabled', false,
                                'MIGProfile', '',
                                'Labels', JSON_OBJECT(),
                                'Annotations', JSON_OBJECT(),
                                'Allocated', false,
                                'Reserved', false
                        ),
                        JSON_OBJECT(
                                'ID', 'gpu-a1',
                                'UUID', 'GPU-A1-UUID',
                                'NodeName', 'node-a',
                                'Index', 1,
                                'Model', 'A100',
                                'Vendor', 'nvidia',
                                'Type', 'full',
                                'MemoryMiB', 81920,
                                'FreeMemoryMiB', 81920,
                                'Healthy', true,
                                'Health', 'healthy',
                                'MIGEnabled', false,
                                'MIGProfile', '',
                                'Labels', JSON_OBJECT(),
                                'Annotations', JSON_OBJECT(),
                                'Allocated', false,
                                'Reserved', false
                        )
                        ),
                'migs', JSON_ARRAY()
        ),
        UTC_TIMESTAMP()
    ),
    (
        'v1',
        'node-b',
        'ready',
        1,
        JSON_OBJECT('gpu.vendor', 'nvidia', 'zone', 'az-b'),
        JSON_OBJECT(),
        JSON_OBJECT(
                'labels', JSON_OBJECT('gpu.vendor', 'nvidia', 'zone', 'az-b'),
                'annotations', JSON_OBJECT(),
                'topology', JSON_OBJECT(
                        'NodeName', 'node-b',
                        'Links', JSON_ARRAY()
                            ),
                'gpus', JSON_ARRAY(
                        JSON_OBJECT(
                                'ID', 'gpu-b0',
                                'UUID', 'GPU-B0-UUID',
                                'NodeName', 'node-b',
                                'Index', 0,
                                'Model', '4090',
                                'Vendor', 'nvidia',
                                'Type', 'full',
                                'MemoryMiB', 24576,
                                'FreeMemoryMiB', 24576,
                                'Healthy', true,
                                'Health', 'healthy',
                                'MIGEnabled', false,
                                'MIGProfile', '',
                                'Labels', JSON_OBJECT(),
                                'Annotations', JSON_OBJECT(),
                                'Allocated', false,
                                'Reserved', false
                        )
                        ),
                'migs', JSON_ARRAY()
        ),
        UTC_TIMESTAMP()
    );

-- gpu_devices: bind to latest two node_snapshots rows for version v1
INSERT INTO gpu_devices (
    id,
    snapshot_id,
    node_name,
    uuid,
    gpu_index,
    model,
    vendor,
    type,
    memory_mib,
    free_memory_mib,
    healthy,
    health,
    mig_enabled,
    mig_profile,
    labels_json,
    annotations_json,
    allocated,
    reserved,
    created_at
)
SELECT 'gpu-a0', ns.id, 'node-a', 'GPU-A0-UUID', 0, 'A100', 'nvidia', 'full', 81920, 81920, 1, 'healthy', 0, '',
       JSON_OBJECT(), JSON_OBJECT(), 0, 0, UTC_TIMESTAMP()
FROM node_snapshots ns
WHERE ns.version = 'v1' AND ns.node_name = 'node-a'
ORDER BY ns.id DESC
    LIMIT 1
ON DUPLICATE KEY UPDATE
                     snapshot_id = VALUES(snapshot_id),
                     free_memory_mib = VALUES(free_memory_mib),
                     healthy = VALUES(healthy),
                     health = VALUES(health),
                     allocated = VALUES(allocated),
                     reserved = VALUES(reserved);

INSERT INTO gpu_devices (
    id,
    snapshot_id,
    node_name,
    uuid,
    gpu_index,
    model,
    vendor,
    type,
    memory_mib,
    free_memory_mib,
    healthy,
    health,
    mig_enabled,
    mig_profile,
    labels_json,
    annotations_json,
    allocated,
    reserved,
    created_at
)
SELECT 'gpu-a1', ns.id, 'node-a', 'GPU-A1-UUID', 1, 'A100', 'nvidia', 'full', 81920, 81920, 1, 'healthy', 0, '',
       JSON_OBJECT(), JSON_OBJECT(), 0, 0, UTC_TIMESTAMP()
FROM node_snapshots ns
WHERE ns.version = 'v1' AND ns.node_name = 'node-a'
ORDER BY ns.id DESC
    LIMIT 1
ON DUPLICATE KEY UPDATE
                     snapshot_id = VALUES(snapshot_id),
                     free_memory_mib = VALUES(free_memory_mib),
                     healthy = VALUES(healthy),
                     health = VALUES(health),
                     allocated = VALUES(allocated),
                     reserved = VALUES(reserved);

INSERT INTO gpu_devices (
    id,
    snapshot_id,
    node_name,
    uuid,
    gpu_index,
    model,
    vendor,
    type,
    memory_mib,
    free_memory_mib,
    healthy,
    health,
    mig_enabled,
    mig_profile,
    labels_json,
    annotations_json,
    allocated,
    reserved,
    created_at
)
SELECT 'gpu-b0', ns.id, 'node-b', 'GPU-B0-UUID', 0, '4090', 'nvidia', 'full', 24576, 24576, 1, 'healthy', 0, '',
       JSON_OBJECT(), JSON_OBJECT(), 0, 0, UTC_TIMESTAMP()
FROM node_snapshots ns
WHERE ns.version = 'v1' AND ns.node_name = 'node-b'
ORDER BY ns.id DESC
    LIMIT 1
ON DUPLICATE KEY UPDATE
                     snapshot_id = VALUES(snapshot_id),
                     free_memory_mib = VALUES(free_memory_mib),
                     healthy = VALUES(healthy),
                     health = VALUES(health),
                     allocated = VALUES(allocated),
                     reserved = VALUES(reserved);

SET FOREIGN_KEY_CHECKS = 1;