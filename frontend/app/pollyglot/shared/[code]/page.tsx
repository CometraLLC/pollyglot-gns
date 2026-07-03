import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { SharedDeckPage } from '@/src/presentation/components/pages/shared-deck-page'

export default async function Page({ params }: { params: Promise<{ code: string }> }) {
	const { code } = await params
	return (
		<ProtectedRoute>
			<MainLayout>
				<SharedDeckPage code={code} />
			</MainLayout>
		</ProtectedRoute>
	)
}
