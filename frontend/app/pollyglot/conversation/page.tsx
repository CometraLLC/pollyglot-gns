import { MainLayout } from '@/src/presentation/components/layout/main-layout'
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route'
import { ConversationPage } from '@/src/presentation/components/pages/conversation-page'

export default function Page() {
	return (
		<ProtectedRoute>
			<MainLayout>
				<ConversationPage />
			</MainLayout>
		</ProtectedRoute>
	)
}
