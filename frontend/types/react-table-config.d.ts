// TypeScript Version: 3.5
// reflects react-table@7.0.0-beta.23

// note that if you are using this you might well want to investigate the use of
// patch-package to delete the index.d.ts file in the react-table package

declare module 'react-table' {
  import {
    ComponentType,
    DependencyList,
    EffectCallback,
    MouseEvent,
    ReactElement,
    ReactNode,
  } from 'react'

  /**
   * The empty definitions of below provides a base definition for the parts used by useTable, that can then be extended in the users code.
   *
   * @example
   *  export interface TableOptions<D extends Record<string, unknown> = {}}>
   *    extends
   *      UseExpandedOptions<D>,
   *      UseFiltersOptions<D> {}
   */
  export interface TableOptions<D extends Record<string, unknown>>
   extends UseExpandedOptions<D>,
      UseTableOptions<D>,
      UsePaginationOptions<D>,
      UseSortByOptions<D> {}
  // export interface TableOptions<D extends Record<string, unknown>> extends UseTableOptions<D> {}

  export interface TableInstance<D extends Record<string, unknown>>
    extends UseTableInstanceProps<D>,
      UsePaginationInstanceProps<D> {}

  // eslint-disable-next-line @typescript-eslint/no-empty-interface
  export interface TableState<D extends Record<string, unknown> = Record<string, unknown>>
    extends UseColumnOrderState<D>,
      UseExpandedState<D>,
      UseFiltersState<D>,
      UseGroupByState<D>,
      UsePaginationState<D>,
      UseRowSelectState<D>,
      UseSortByState<D> {
    rowCount: number
  } /* tslint:disable-line no-empty-interface */

  export interface Hooks<D extends Record<string, unknown> = Record<string, unknown>> extends UseTableHooks<D> {}

  export interface Cell<D extends Record<string, unknown> = Record<string, unknown>> extends UseTableCellProps<D> {
    getCellProps(): (propGetter?: CellPropGetter<D>) => TableCellProps;
    render: (type: 'Cell' | string, userProps?: object) => ReactNode;
  }

  export interface HeaderGroup<D extends object = {}> extends ColumnInstance<D>, UseTableHeaderGroupProps<D> {}

  export interface UseTableHeaderGroupProps<D extends Record<string, unknown>> {
    headers: Array<HeaderGroup<D>>;
    getHeaderGroupProps: (propGetter?: HeaderGroupPropGetter<D>) => TableHeaderProps;
    getFooterGroupProps: (propGetter?: FooterGroupPropGetter<D>) => TableFooterProps;
    totalHeaderCount: number; // not documented
  }

export interface UseTableColumnProps<D extends Record<string, unknown>> {
    id: IdType<D>;
    columns?: Array<ColumnInstance<D>>;
    isVisible: boolean;
    render: (type: 'Header' | 'Footer' | string, props?: object) => ReactNode;
    totalLeft: number;
    totalWidth: number;
    getHeaderProps: (propGetter?: HeaderPropGetter<D>) => TableHeaderProps;
    getFooterProps: (propGetter?: FooterPropGetter<D>) => TableFooterProps;
    toggleHidden: (value?: boolean) => void;
    parent?: ColumnInstance<D>; // not documented
    getToggleHiddenProps: (userProps?: any) => any;
    depth: number; // not documented
    placeholderOf?: ColumnInstance;
  }

  export interface Column<D extends Record<string, unknown> = Record<string, unknown>>
    extends UseTableColumnOptions<D> {}

  export interface ColumnInstance<D extends Record<string, unknown> = Record<string, unknown>>
    extends Omit<Column<D>, 'id'>,
      UseTableColumnProps<D> {}

  export interface Row<D extends Record<string, unknown> = Record<string, unknown>> extends UseTableRowProps<D> {}

  /* #region useTable */
  export function useTable<D extends Record<string,unknown> = Record<string, unknown>>(
    options: TableOptions<D>,
    ...plugins: Array<PluginHook<D>>
  ): TableInstance<D>

  /**
   * NOTE: To use custom options, use "Interface Merging" to add the custom options
   */
  export type UseTableOptions<D extends Record<string,unknown>> = {
    columns: Array<Column<D>>
    data: Array<D>
  } & Partial<{
    initialState: Partial<TableState<D>>
    reducer: (
      newState: TableState<D>,
      action: string,
      prevState: TableState<D>
    ) => TableState<D>
    defaultColumn: Partial<Column<D>>
    initialRowStateKey: IdType<D>
    getSubRows: (originalRow: D, relativeIndex: number) => Array<D>
    getRowId: (originalRow: D, relativeIndex: number) => IdType<D>
    debug: boolean
  }>

  export interface UseTableInstanceProps<D extends Record<string,unknown>> {
    columns: Array<ColumnInstance<D>>
    flatColumns: Array<ColumnInstance<D>>
    headerGroups: Array<HeaderGroup<D>>
    headers: Array<ColumnInstance<D>>
    flatHeaders: Array<ColumnInstance<D>>
    rows: Array<Row<D>>
    getTableProps: (props?: Record<string,unknown>) => Record<string,unknown>
    getTableBodyProps: (props?: Record<string,unknown>) => Record<string,unknown>
    prepareRow: (row: Row<D>) => void
    rowPaths: string[]
    flatRows: Array<Row<D>>
    state: TableState<D>
    dispatch: TableDispatch<D, TableAction>
    totalColumnsWidth: number
  }

  export interface UseTableRowProps<D extends Record<string, unknown>> {
    cells: Array<Cell<D>>
    values: Record<IdType<D>, CellValue>
    getRowProps: (props?: Record<string, unknown>) => Record<string, unknown>
    index: number
    original: D
    path: Array<IdType<D>>
    subRows: Array<Row<D>>
  }

  // NOTE: At least one of (id | accessor | Header as string) required
  export interface Accessor<D extends Record<string, unknown>> {
    accessor?:
      | IdType<D>
      | ((
          originalRow: D,
          index: number,
          sub: {
            subRows: D[]
            depth: number
            data: D[]
          }
        ) => CellValue)
    id?: IdType<D>
  }

  /* #endregion */

  // Plugins

  /* #region useColumnOrder */
  export function useColumnOrder<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useColumnOrder {
    const pluginName = 'useColumnOrder'
  }

  export interface UseColumnOrderState<D extends Record<string, unknown>> {
    columnOrder: Array<IdType<D>>
  }

  export interface UseColumnOrderInstanceProps<D extends Record<string, unknown>> {
    setColumnOrder: (
      updater: (columnOrder: Array<IdType<D>>) => Array<IdType<D>>
    ) => void
  }

  /* #endregion */

  /* #region useExpanded */
  export function useExpanded<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useExpanded {
    const pluginName = 'useExpanded'
  }

  export type UseExpandedOptions<D extends Record<string, unknown>> = Partial<{
    manualExpandedKey: IdType<D>
    paginateExpandedRows: boolean
    getResetExpandedDeps: (instance: TableInstance<D>) => Array<any>
  }>

  export interface UseExpandedHooks<D extends Record<string, unknown>> {
    getExpandedToggleProps: Array<
      (row: Row<D>, instance: TableInstance<D>) => Record<string, unknown>
    >
  }

  export interface UseExpandedState<D extends Record<string, unknown>> {
    expanded: Array<IdType<D>>
  }

  export interface UseExpandedInstanceProps<D extends Record<string, unknown>> {
    rows: Array<Row<D>>
    toggleExpandedByPath: (path: Array<IdType<D>>, isExpanded: boolean) => void
    expandedDepth: number
  }

  export interface UseExpandedRowProps<D extends Record<string, unknown>> {
    isExpanded: boolean
    canExpand: boolean
    subRows: Array<Row<D>>
    toggleExpanded: (isExpanded?: boolean) => void
    getExpandedToggleProps: (props?: Record<string, unknown>) => Record<string, unknown>
  }

  /* #endregion */

  /* #region useFilters */
  export function useFilters<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useFilters {
    const pluginName = 'useFilters'
  }

  export type UseFiltersOptions<D extends Record<string, unknown>> = Partial<{
    manualFilters: boolean
    disableFilters: boolean
    defaultCanFilter: boolean
    filterTypes: Filters<D>
    getResetFiltersDeps: (instance: TableInstance<D>) => Array<any>
  }>

  export interface UseFiltersState<D extends Record<string, unknown>> {
    filters: Filters<D>
  }

  export type UseFiltersColumnOptions<D extends Record<string, unknown>> = Partial<{
    Filter: Renderer<FilterProps<D>>
    disableFilters: boolean
    defaultCanFilter: boolean
    filter: FilterType<D> | DefaultFilterTypes | keyof Filters<D>
  }>

  export interface UseFiltersInstanceProps<D extends Record<string, unknown>> {
    rows: Array<Row<D>>
    preFilteredRows: Array<Row<D>>
    setFilter: (
      columnId: IdType<D>,
      updater: ((filterValue: FilterValue) => FilterValue) | FilterValue
    ) => void
    setAllFilters: (
      updater: Filters<D> | ((filters: Filters<D>) => Filters<D>)
    ) => void
  }

  export interface UseFiltersColumnProps<D extends Record<string, unknown>> {
    canFilter: boolean
    setFilter: (
      updater: ((filterValue: FilterValue) => FilterValue) | FilterValue
    ) => void
    filterValue: FilterValue
    preFilteredRows: Array<Row<D>>
    filteredRows: Array<Row<D>>
  }

  export type FilterProps<D extends Record<string, unknown>> = HeaderProps<D>
  export type FilterValue = any
  export type Filters<D extends Record<string, unknown>> = Record<IdType<D>, FilterValue>

  export type DefaultFilterTypes =
    | 'text'
    | 'exactText'
    | 'exactTextCase'
    | 'includes'
    | 'includesAll'
    | 'exact'
    | 'equals'
    | 'between'

  export interface FilterType<D extends Record<string, unknown>> {
    (
      rows: Array<Row<D>>,
      columnId: IdType<D>,
      filterValue: FilterValue,
      column: ColumnInstance<D>
    ): Array<Row<D>>

    autoRemove?: (filterValue: FilterValue) => boolean
  }

  /* #endregion */

  /* #region useGroupBy */
  export function useGroupBy<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useGroupBy {
    const pluginName = 'useGroupBy'
  }

  export type UseGroupByOptions<D extends Record<string, unknown>> = Partial<{
    manualGroupBy: boolean
    disableGroupBy: boolean
    defaultCanGroupBy: boolean
    aggregations: Record<string, AggregatorFn<D>>
    groupByFn: (
      rows: Array<Row<D>>,
      columnId: IdType<D>
    ) => Record<string, Row<D>>
    getResetGroupByDeps: (instance: TableInstance<D>) => Array<any>
  }>

  export interface UseGroupByHooks<D extends Record<string, unknown>> {
    getGroupByToggleProps: Array<
      (header: HeaderGroup<D>, instance: TableInstance<D>) => Record<string, unknown>
    >
  }

  export interface UseGroupByState<D extends Record<string, unknown>> {
    groupBy: Array<IdType<D>>
  }

  export type UseGroupByColumnOptions<D extends Record<string, unknown>> = Partial<{
    aggregate: Aggregator<D> | Array<Aggregator<D>>
    Aggregated: Renderer<CellProps<D>>
    disableGroupBy: boolean
    defaultCanGroupBy: boolean
    groupByBoundary: boolean
  }>

  export interface UseGroupByInstanceProps<D extends Record<string, unknown>> {
    rows: Array<Row<D>>
    preGroupedRows: Array<Row<D>>
    toggleGroupBy: (columnId: IdType<D>, toggle: boolean) => void
  }

  export interface UseGroupByColumnProps<D extends Record<string, unknown>> {
    canGroupBy: boolean
    isGrouped: boolean
    groupedIndex: number
    toggleGroupBy: () => void
    getGroupByToggleProps: (props?: Record<string, unknown>) => Record<string, unknown>
  }

  export interface UseGroupByRowProps<D extends Record<string, unknown>> {
    isAggregated: boolean
    groupByID: IdType<D>
    groupByVal: string
    values: Record<IdType<D>, AggregatedValue>
    subRows: Array<Row<D>>
    depth: number
    path: Array<IdType<D>>
    index: number
  }

  export interface UseGroupByCellProps<D extends Record<string, unknown>> {
    isGrouped: boolean
    isRepeatedValue: boolean
    isAggregated: boolean
  }

  export type DefaultAggregators =
    | 'sum'
    | 'average'
    | 'median'
    | 'uniqueCount'
    | 'count'

  export type AggregatorFn<D extends Record<string, unknown>> = (
    columnValues: CellValue[],
    rows: Array<Row<D>>,
    isAggregated: boolean
  ) => AggregatedValue
  export type Aggregator<D extends Record<string, unknown>> =
    | AggregatorFn<D>
    | DefaultAggregators
    | string
  export type AggregatedValue = any
  /* #endregion */

  /* #region usePagination */
  export function usePagination<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace usePagination {
    const pluginName = 'usePagination'
  }

  export type UsePaginationOptions<D extends Record<string, unknown>> = Partial<{
    pageCount: number
    manualPagination: boolean
    getResetPageDeps: (instance: TableInstance<D>) => Array<any>
    paginateExpandedRows: boolean
  }>

  export interface UsePaginationState<D extends Record<string, unknown>> {
    pageSize: number
    pageIndex: number
  }

  export interface UsePaginationInstanceProps<D extends Record<string, unknown>> {
    page: Array<Row<D>>
    pageCount: number
    pageOptions: number[]
    canPreviousPage: boolean
    canNextPage: boolean
    gotoPage: (updater: ((pageIndex: number) => number) | number) => void
    previousPage: () => void
    nextPage: () => void
    setPageSize: (pageSize: number) => void
    pageIndex: number
    pageSize: number
  }

  /* #endregion */

  /* #region useRowSelect */
  export function useRowSelect<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useRowSelect {
    const pluginName = 'useRowSelect'
  }

  export type UseRowSelectOptions<D extends Record<string, unknown>> = Partial<{
    manualRowSelectedKey: IdType<D>
    getResetSelectedRowPathsDeps: (instance: TableInstance<D>) => Array<any>
  }>

  export interface UseRowSelectHooks<D extends Record<string, unknown>> {
    getToggleRowSelectedProps: Array<
      (row: Row<D>, instance: TableInstance<D>) => Record<string, unknown>
    >
    getToggleAllRowsSelectedProps: Array<(instance: TableInstance<D>) => Record<string, unknown>>
  }

  export interface UseRowSelectState<D extends Record<string, unknown>> {
    selectedRowPaths: Array<IdType<D>>
  }

  export interface UseRowSelectInstanceProps<D extends Record<string, unknown>> {
    toggleRowSelected: (rowPath: IdType<D>, set?: boolean) => void
    toggleRowSelectedAll: (set?: boolean) => void
    getToggleAllRowsSelectedProps: (props?: Record<string, unknown>) => Record<string, unknown>
    isAllRowsSelected: boolean
    selectedFlatRows: Array<Row<D>>
  }

  export interface UseRowSelectRowProps<D extends Record<string, unknown>> {
    isSelected: boolean
    toggleRowSelected: (set?: boolean) => void
    getToggleRowSelectedProps: (props?: Record<string, unknown>) => Record<string, unknown>
  }

  /* #endregion */

  /* #region useRowState */
  export function useRowState<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useRowState {
    const pluginName = 'useRowState'
  }

  export type UseRowStateOptions<D extends Record<string, unknown>> = Partial<{
    initialRowStateAccessor: (row: Row<D>) => UseRowStateLocalState<D>
  }>

  export interface UseRowStateState<D extends Record<string, unknown>> {
    rowState: Partial<{
      cellState: UseRowStateLocalState<D>
      rowState: UseRowStateLocalState<D>
    }>
  }

  export interface UseRowStateInstanceProps<D extends Record<string, unknown>> {
    setRowState: (rowPath: string[], updater: UseRowUpdater) => void // Purposely not exposing action
    setCellState: (
      rowPath: string[],
      columnId: IdType<D>,
      updater: UseRowUpdater
    ) => void
  }

  export interface UseRowStateRowProps<D extends Record<string, unknown>> {
    state: UseRowStateLocalState<D>
    setState: (updater: UseRowUpdater) => void
  }

  export interface UseRowStateCellProps<D extends Record<string, unknown>> {
    state: UseRowStateLocalState<D>
    setState: (updater: UseRowUpdater) => void
  }

  export type UseRowUpdater<T = unknown> = T | ((prev: T) => T)
  export type UseRowStateLocalState<D extends Record<string, unknown>, T = unknown> = Record<
    IdType<D>,
    T
  >
  /* #endregion */

  /* #region useSortBy */
  export function useSortBy<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useSortBy {
    const pluginName = 'useSortBy'
  }

  export type UseSortByOptions<D extends Record<string, unknown>> = Partial<{
    manualSorting: boolean
    disableSortBy: boolean
    defaultCanSort: boolean
    disableMultiSort: boolean
    isMultiSortEvent: (e: MouseEvent) => boolean
    maxMultiSortColCount: number
    disableSortRemove: boolean
    disabledMultiRemove: boolean
    orderByFn: (
      rows: Array<Row<D>>,
      sortFns: Array<SortByFn<D>>,
      directions: boolean[]
    ) => Array<Row<D>>
    sortTypes: Record<string, SortByFn<D>>
    getResetSortByDeps: (instance: TableInstance<D>) => Array<any>
  }>

  export interface UseSortByHooks<D extends Record<string, unknown>> {
    getSortByToggleProps: Array<
      (column: Column<D>, instance: TableInstance<D>) => Record<string, unknown>
    >
  }

  export interface UseSortByState<D extends Record<string, unknown>> {
    sortBy: Array<SortingRule<D>>
  }

  export type UseSortByColumnOptions<D extends Record<string, unknown>> = Partial<{
    defaultCanSort: boolean
    disableSortBy: boolean
    sortDescFirst: boolean
    sortInverted: boolean
    sortType: SortByFn<D> | DefaultSortTypes | string
  }>

  export interface UseSortByInstanceProps<D extends Record<string, unknown>> {
    rows: Array<Row<D>>
    preSortedRows: Array<Row<D>>
    toggleSortBy: (
      columnId: IdType<D>,
      descending: boolean,
      isMulti: boolean
    ) => void
  }

  export interface UseSortByColumnProps<D extends Record<string, unknown>> {
    canSort: boolean
    toggleSortBy: (descending: boolean, multi: boolean) => void
    getSortByToggleProps: (props?: Record<string, unknown>) => Record<string, unknown>
    clearSorting: () => void
    isSorted: boolean
    sortedIndex: number
    isSortedDesc: boolean | undefined
  }

  export type SortByFn<D extends Record<string, unknown>> = (
    rowA: Row<D>,
    rowB: Row<D>,
    columnId: IdType<D>
  ) => 0 | 1 | -1

  export type DefaultSortTypes = 'alphanumeric' | 'datetime' | 'basic'

  export interface SortingRule<D> {
    id: IdType<D>
    desc?: boolean
  }

  /* #endregion */

  /* #region useAbsoluteLayout */
  export function useAbsoluteLayout<D extends Record<string, unknown> = {}>(
    hooks: Hooks<D>
  ): void

  export namespace useAbsoluteLayout {
    const pluginName = 'useAbsoluteLayout'
  }
  /* #endregion */

  /* #region useBlockLayout */
  export function useBlockLayout<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useBlockLayout {
    const pluginName = 'useBlockLayout'
  }
  /* #endregion */

  /* #region useResizeColumns */
  export function useResizeColumns<D extends Record<string, unknown> = {}>(hooks: Hooks<D>): void

  export namespace useResizeColumns {
    const pluginName = 'useResizeColumns'
  }

  export interface UseResizeColumnsOptions<D extends Record<string, unknown>> {
    disableResizing?: boolean
  }

  export interface UseResizeColumnsColumnOptions<D extends Record<string, unknown>> {
    disableResizing?: boolean
  }

  export interface UseResizeColumnsHeaderProps<D extends Record<string, unknown>> {
    getResizerProps: (props?: Record<string, unknown>) => Record<string, unknown>
    canResize: boolean
    isResizing: boolean
  }

  /* #endregion */

  type ElementDimensions = {
    left: number
    width: number
    outerWidth: number
    marginLeft: number
    marginRight: number
    paddingLeft: number
    paddingRight: number
    scrollWidth: number
  }


  // Additional API
  export const actions: Record<string, string>
  export const defaultColumn: Partial<Column>
  type TableReducer<S, A> = (prevState: S, action: A) => S | undefined
  export const reducerHandlers: Record<string, TableReducer<TableState, any>>

  // Helpers
  export type StringKey<D> = Extract<keyof D, string>
  export type IdType<D> = StringKey<D> | string
  export type CellValue = any

  export type Renderer<Props> = ComponentType<Props> | ReactNode

  export interface PluginHook<D extends Record<string, unknown>> {
    (hooks: Hooks<D>): void
    pluginName: string
  }

  export type TableAction = {
    type: keyof typeof actions | string
  }

  export type TableDispatch<D extends Record<string, unknown>, A extends TableAction> = (
    action: A
  ) => void
}
