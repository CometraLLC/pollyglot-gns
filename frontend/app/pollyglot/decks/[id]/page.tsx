import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { DeckDetailPage } from '@/src/presentation/components/pages/deck-detail-page'

export default async function Page({ params }: { params: Promise<{ id: string }> }) {
	const { id } = await params
	return (
		<ProtectedRoute>
			<MainLayout>
				<DeckDetailPage deckId={id} />
			</MainLayout>
		</ProtectedRoute>
	)
}
