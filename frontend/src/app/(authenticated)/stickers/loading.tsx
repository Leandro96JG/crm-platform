import PageHeaderSkeleton from '@/app/components/common/skeletons/page-header-skeleton';
import TableSkeleton from '@/app/components/common/skeletons/table-skeleton';

export default function Loading() {
  return (
    <>
      <PageHeaderSkeleton />
      <main className="px-8 pb-10">
        <TableSkeleton
          headers={['Plancha', 'Categoría', 'Subcategoría', 'Estado']}
        />
      </main>
    </>
  );
}
