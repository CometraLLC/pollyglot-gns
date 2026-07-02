import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { DecksPage } from '@/src/presentation/components/pages/decks-page'

export default function Page() {
	return (
		<ProtectedRoute>
			<MainLayout>
				<DecksPage />
			</MainLayout>
		</ProtectedRoute>
	)
}
