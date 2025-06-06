class UserPolicy
  attr_reader :current_user, :record

  def initialize(current_user, record)
    @current_user = current_user
    @record = record
  end

  def create?
    return record.role == 'member' if current_user.nil? # Guest signup

    current_user.admin? && record.role == 'librarian'
  end

  def update?
    return false unless current_user.present?

    # Admin can update librarians
    return true if current_user.admin? && record.librarian?

    # Members can update themselves
    current_user == record && current_user.member?
  end

  def destroy?
    return false unless current_user.present?

    # Admin can delete librarians
    return true if current_user.admin? && record.librarian?

    # Members can delete themselves
    current_user == record && current_user.member?
  end
end
