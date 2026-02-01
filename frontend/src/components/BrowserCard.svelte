<script>
  export let instance;
  export let onStart;
  export let onStop;
  export let onEdit;
  export let onDuplicate;
  export let onDelete;
  export let isMiniMode = false;

  $: statusColor = instance.running ? 'var(--success-color)' : 'var(--text-muted)';
  $: statusText = instance.running ? '运行中' : '已停止';
</script>

<div class="card" class:mini={isMiniMode}>
  <div class="card-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: {isMiniMode ? '0' : '1rem'};">
    <div style="flex: 1; min-width: 0; margin-right: 0.5rem;">
      <h3 style="font-size: 1.125rem; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;" title={instance.name}>{instance.name}</h3>
      {#if !isMiniMode}
      <div style="display: flex; align-items: center; gap: 0.5rem; margin-top: 0.25rem;">
        <span style="width: 8px; height: 8px; border-radius: 50%; background-color: {statusColor}"></span>
        <span style="font-size: 0.75rem; color: var(--text-muted); font-weight: 500;">{statusText}</span>
      </div>
      {/if}
    </div>
    
    {#if !isMiniMode}
    <div class="actions" style="display: flex; gap: 0.5rem;">
      <button class="btn-icon" on:click={() => onEdit(instance)} title="编辑">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
      </button>
      <button class="btn-icon" on:click={() => onDuplicate(instance)} title="复制配置">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
      </button>
      <button class="btn-icon" on:click={() => onDelete(instance.id)} title="删除" style="color: var(--danger-color);">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
      </button>
    </div>
    {:else}
      <!-- Mini Mode Action: Only Start/Stop Button -->
      {#if !instance.running}
        <button class="btn btn-primary btn-sm" style="padding: 0.25rem 0.75rem;" on:click={() => onStart(instance.id)}>
          启动
        </button>
      {:else}
        <button class="btn btn-sm" style="padding: 0.25rem 0.75rem; border-color: var(--danger-color); color: var(--danger-color);" on:click={() => onStop(instance.id)}>
          停止
        </button>
      {/if}
    {/if}
  </div>

  {#if !isMiniMode}
  <div class="card-body" style="margin-bottom: 1.5rem;">
    <div class="info-item">
      <span class="label">启动路径:</span>
      <span class="value" title={instance.path}>{instance.path}</span>
    </div>
    <div class="info-item">
      <span class="label">数据目录:</span>
      <span class="value" title={instance.userDataDir}>{instance.userDataDir || '默认'}</span>
    </div>
    <div class="info-item">
        <span class="label">参数:</span>
        <span class="value">{instance.args.join(' ') || '无'}</span>
      </div>
  </div>

  <div class="card-footer" style="display: flex; gap: 0.75rem;">
    {#if !instance.running}
      <button class="btn btn-primary" style="flex: 1;" on:click={() => onStart(instance.id)}>
        启动浏览器
      </button>
    {:else}
      <button class="btn" style="flex: 1; border-color: var(--danger-color); color: var(--danger-color);" on:click={() => onStop(instance.id)}>
        停止
      </button>
    {/if}
  </div>
  {/if}
</div>

<style>
  .btn-icon {
    background: transparent;
    border: none;
    padding: 0.25rem;
    cursor: pointer;
    color: var(--text-muted);
    border-radius: 0.25rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .btn-icon:hover {
    background: #f1f5f9;
    color: var(--text-main);
  }
  .info-item {
    display: flex;
    flex-direction: column;
    margin-bottom: 0.5rem;
  }
  .label {
    font-size: 0.75rem;
    color: var(--text-muted);
    margin-bottom: 0.125rem;
  }
  .value {
    font-size: 0.8125rem;
    color: var(--text-main);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
